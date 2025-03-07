package unify

import (
	"errors"
	"hw4/disjointset"
	"hw4/term"
)

// ErrUnifier is returned when two terms cannot be unified.
var ErrUnifier = errors.New("unifier error")

// UnifyResult is the result of unification. For example, for a variable term
// `s`, UnifyResult[s] is the term that s is unified with.
type UnifyResult map[*term.Term]*term.Term

// Unifier is the interface for the term unifier.
// Do not change the definition of this interface.
type Unifier interface {
	Unify(*term.Term, *term.Term) (UnifyResult, error)
}

// UnifierImpl implements the Unifier interface using a disjoint-set
// to manage equivalences between variables and a substitution map to record
// non-variable bindings.
type UnifierImpl struct {
	ds     disjointset.DisjointSet   // union–find for variable equivalences
	varID  map[*term.Term]int        // assigns a unique integer id to each variable term
	repVar map[int]*term.Term        // stores a canonical variable for each DS set (by rep id)
	subst  map[int]*term.Term        // maps a DS representative id to a binding term (if any)
	nextID int                       // next unique id to assign
}

// NewUnifier creates an instance of UnifierImpl.
func NewUnifier() Unifier {
	return &UnifierImpl{
		ds:     disjointset.NewDisjointSet(),
		varID:  make(map[*term.Term]int),
		repVar: make(map[int]*term.Term),
		subst:  make(map[int]*term.Term),
		nextID: 0,
	}
}

// getID returns a unique id for the variable term v.
// If v does not already have an id, one is assigned and v is recorded as the canonical variable.
func (u *UnifierImpl) getID(v *term.Term) int {
	if id, ok := u.varID[v]; ok {
		return id
	}
	id := u.nextID
	u.nextID++
	u.varID[v] = id
	u.repVar[id] = v
	return id
}

// prune returns the current representative (or binding) of term t.
// For a variable, if a binding exists in subst, we recursively prune that binding;
// otherwise we return the canonical variable from repVar.
// For compound terms, we prune each argument.
func (u *UnifierImpl) prune(t *term.Term) *term.Term {
	if t == nil {
	}
	if t.Typ == term.TermVariable {
		id := u.getID(t)
		rep := u.ds.FindSet(id)
		if binding, ok := u.subst[rep]; ok {
			pruned := u.prune(binding)
			u.subst[rep] = pruned
			return pruned
		}
		return u.repVar[rep]
	}
	if t.Typ == term.TermCompound {
		newArgs := make([]*term.Term, len(t.Args))
		for i, arg := range t.Args {
			newArgs[i] = u.prune(arg)
		}
		return &term.Term{
			Typ:     t.Typ,
			Functor: t.Functor,
			Args:    newArgs,
		}
	}
	// Atoms and numbers are returned as is.
	return t
}

// occursCheck ensures that variable v does not occur in term t.
// This prevents cycles (e.g. unifying X with f(X)).
func (u *UnifierImpl) occursCheck(v *term.Term, t *term.Term) bool {
	t = u.prune(t)
	if t.Typ == term.TermVariable {
		// If the canonical representatives are equal then v occurs in t.
		return u.ds.FindSet(u.getID(v)) == u.ds.FindSet(u.getID(t))
	}
	if t.Typ == term.TermCompound {
		for _, arg := range t.Args {
			if u.occursCheck(v, arg) {
				return true
			}
		}
	}
	return false
}

// equalTerms compares two terms for equality.
// (For variables, since the parser interns simple terms, pointer equality is sufficient.)
func (u *UnifierImpl) equalTerms(s, t *term.Term) bool {
	if s.Typ != t.Typ {
		return false
	}
	if s.Typ == term.TermVariable {
		return s == t
	}
	if s.Typ != term.TermCompound {
		return s.Literal == t.Literal
	}
	// For compound terms, compare functor and arguments.
	if s.Functor == nil || t.Functor == nil || s.Functor.Literal != t.Functor.Literal {
		return false
	}
	if len(s.Args) != len(t.Args) {
		return false
	}
	for i := 0; i < len(s.Args); i++ {
		if !u.equalTerms(s.Args[i], t.Args[i]) {
			return false
		}
	}
	return true
}

// unify is a recursive helper that unifies terms s and t.
// It uses the disjoint-set to merge free variables and the subst map for bindings.
// The algorithm performs an occurs check when a variable is about to be bound.
func (u *UnifierImpl) unify(s *term.Term, t *term.Term) error {
	s = u.prune(s)
	t = u.prune(t)
	if u.equalTerms(s, t) {
		return nil
	}
	// If s is a variable...
	if s.Typ == term.TermVariable {
		// If both are variables, merge their equivalence classes.
		if t.Typ == term.TermVariable {
			idS := u.getID(s)
			idT := u.getID(t)
			repS := u.ds.FindSet(idS)
			repT := u.ds.FindSet(idT)
			if repS != repT {
				newRep := u.ds.UnionSet(idS, idT)
					// No bindings exist; update the canonical variable.
					if idS < idT {
						u.repVar[newRep] = s
					} else {
						u.repVar[newRep] = t
					}
			}
			return nil
		}
		// s is variable and t is not.
		if u.occursCheck(s, t) {
			return ErrUnifier
		}
		idS := u.getID(s)
		rep := u.ds.FindSet(idS)
		u.subst[rep] = t
		return nil
	}
	// Similarly, if t is variable.
	if t.Typ == term.TermVariable {
		if u.occursCheck(t, s) {
			return ErrUnifier
		}
		idT := u.getID(t)
		rep := u.ds.FindSet(idT)
		u.subst[rep] = s
		return nil
	}
	// Both s and t are non-variable.
	// If both are atoms or numbers, compare their literals.
	if s.Typ != term.TermCompound && t.Typ != term.TermCompound {
		if s.Literal == t.Literal {
		}
		return ErrUnifier
	}
	// If one is a compound term and the other is not, they cannot be unified.
	if s.Typ != t.Typ {
		return ErrUnifier
	}
	// Both are compound: check that the functors match.
	if s.Functor == nil || t.Functor == nil || s.Functor.Literal != t.Functor.Literal {
		return ErrUnifier
	}
	// Check that the number of arguments is the same.
	if len(s.Args) != len(t.Args) {
		return ErrUnifier
	}
	// Recursively unify each pair of corresponding arguments.
	for i := 0; i < len(s.Args); i++ {
		if err := u.unify(s.Args[i], t.Args[i]); err != nil {
			return err
		}
	}
	return nil
}

// Unify is the method required by the Unifier interface.
// It calls the helper unify method and then builds a substitution mapping
// (the most general unifier) by iterating over all variable ids.
func (u *UnifierImpl) Unify(s, t *term.Term) (UnifyResult, error) {
	if err := u.unify(s, t); err != nil {
		return nil, err
	}
	result := make(UnifyResult)
	// For every variable encountered, determine its binding.
	for v, id := range u.varID {
		rep := u.ds.FindSet(id)
		if binding, ok := u.subst[rep]; ok {
			pruned := u.prune(binding)
			if pruned != v { // only include if the binding is not the variable itself
				result[v] = pruned
			}
		} else {
			// For free variables, choose a canonical representative.
			canonical := u.repVar[rep]
			if canonical != v {
				result[v] = canonical
			}
		}
	}
	return result, nil
}
