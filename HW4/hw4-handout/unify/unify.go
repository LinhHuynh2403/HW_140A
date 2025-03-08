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

// GeneralUnifier implements the Unifier interface using a single disjoint–set.
type GeneralUnifier struct {
	ds     disjointset.DisjointSet // single disjoint-set instance for all terms
	size   map[int]int             // union–by–size (keyed by term index)
	schema map[int]*term.Term      // each set's “canonical” term (variable or non-variable)
	vars   map[int][]*term.Term    // variables in the equivalence class (keyed by representative index)
}

// Global maps for term indexing and cycle detection.
var (
	unifyMap    UnifyResult
	mapToInt    = map[*term.Term]int{}
	// mapToTerm is no longer used for canonical lookup; use unif.schema instead.
	nodeCounter = 0
	visited     = map[int]bool{}
	acyclic     = map[int]bool{}
)

// resetGlobal resets the global maps used for indexing and cycle detection.
func resetGlobal() {
	unifyMap = UnifyResult{}
	mapToInt = map[*term.Term]int{}
	nodeCounter = 0
	visited = map[int]bool{}
	acyclic = map[int]bool{}
}

// NewUnifier creates a new GeneralUnifier.
func NewUnifier() Unifier {
	return &GeneralUnifier{
		ds:     disjointset.NewDisjointSet(),
		size:   make(map[int]int),
		schema: make(map[int]*term.Term),
		vars:   make(map[int][]*term.Term),
	}
}

// Initializer registers a term (if not nil) in the global maps.
func (unif *GeneralUnifier) Initializer(t1 *term.Term, t2 *term.Term) {
	if t1 != nil {
		if _, ok := mapToInt[t1]; !ok {
			mapToInt[t1] = nodeCounter
			unif.size[nodeCounter] = 1
			unif.schema[nodeCounter] = t1
			// Record variable only if t1 is a variable.
			if t1.Typ == term.TermVariable {
				unif.vars[nodeCounter] = []*term.Term{t1}
			} else {
				unif.vars[nodeCounter] = []*term.Term{}
			}
			nodeCounter++
		}
	}
	if t2 != nil {
		if _, ok := mapToInt[t2]; !ok {
			mapToInt[t2] = nodeCounter
			unif.size[nodeCounter] = 1
			unif.schema[nodeCounter] = t2
			if t2.Typ == term.TermVariable {
				unif.vars[nodeCounter] = []*term.Term{t2}
			} else {
				unif.vars[nodeCounter] = []*term.Term{}
			}
			nodeCounter++
		}
	}
}

// occursCheck recursively checks whether varTerm occurs in t.
func occursCheck(varTerm *term.Term, t *term.Term) bool {
	if varTerm == t {
		return true
	}
	if t.Typ == term.TermCompound {
		if occursCheck(varTerm, t.Functor) {
			return true
		}
		for _, arg := range t.Args {
			if occursCheck(varTerm, arg) {
				return true
			}
		}
	}
	return false
}

// Unify attempts to unify two terms.
func (unif *GeneralUnifier) Unify(t1 *term.Term, t2 *term.Term) (UnifyResult, error) {
	resetGlobal()             // Reset global indexing for each unification.
	unif.Initializer(t1, t2)    // Register t1 and t2.
	unifyMap = UnifyResult{}    // Clear previous substitution.
	if err := unif.UnifClosure(t1, t2); err != nil {
		return nil, ErrUnifier
	}
	if err := unif.FindSolution(t1); err != nil {
		return nil, ErrUnifier
	}
	return unifyMap, nil
}

// UnifClosure performs the main unification procedure.
func (unif *GeneralUnifier) UnifClosure(t1 *term.Term, t2 *term.Term) error {
	unif.Initializer(t1, t2)
	idx1 := mapToInt[t1]
	idx2 := mapToInt[t2]
	rep1 := unif.ds.FindSet(idx1)
	rep2 := unif.ds.FindSet(idx2)
	s := unif.schema[rep1]
	t := unif.schema[rep2]

	// If the canonical terms are already identical, nothing to do.
	if s == t {
		return nil
	}

	// Occurs check: if one side is a variable and occurs inside the other.
	if s.Typ == term.TermVariable && t.Typ != term.TermVariable {
		if occursCheck(s, t) {
			return ErrUnifier
		}
	}
	if t.Typ == term.TermVariable && s.Typ != term.TermVariable {
		if occursCheck(t, s) {
			return ErrUnifier
		}
	}

	// If either side is variable, union their equivalence classes.
	if s.Typ == term.TermVariable || t.Typ == term.TermVariable {
		unif.Union(rep1, rep2)
	} else {
		// Both are non-variable: they must be compound with the same functor and arity.
		if s.Typ == term.TermCompound && t.Typ == term.TermCompound {
			if s.Functor != t.Functor || len(s.Args) != len(t.Args) {
				return ErrUnifier
			}
			unif.Union(rep1, rep2)
			// Recursively unify corresponding arguments.
			for i := 0; i < len(s.Args); i++ {
				unif.Initializer(s.Args[i], t.Args[i])
				if err := unif.UnifClosure(s.Args[i], t.Args[i]); err != nil {
					return ErrUnifier
				}
			}
		} else {
			// Mismatched constants (or one constant and one compound) cannot be unified.
			return ErrUnifier
		}
	}
	return nil
}

// Union merges the equivalence classes corresponding to rep1 and rep2.
func (unif *GeneralUnifier) Union(rep1, rep2 int) {
	// Ensure we have the current representatives.
	rep1 = unif.ds.FindSet(rep1)
	rep2 = unif.ds.FindSet(rep2)
	if rep1 == rep2 {
		return
	}
	var newRep, oldRep int
	if unif.size[rep1] >= unif.size[rep2] {
		newRep = rep1
		oldRep = rep2
	} else {
		newRep = rep2
		oldRep = rep1
	}
	unif.size[newRep] += unif.size[oldRep]
	unif.vars[newRep] = append(unif.vars[newRep], unif.vars[oldRep]...)
	// If the canonical term for newRep is a variable, update it to the one from oldRep.
	if unif.schema[newRep].Typ == term.TermVariable {
		unif.schema[newRep] = unif.schema[oldRep]
	}
	// Perform the union in the disjoint-set.
	unif.ds.UnionSet(newRep, oldRep)
}

// FindSolution traverses the unified term structure to build the substitution mapping.
// It also uses a visited/acyclic mechanism to detect cycles.
func (unif *GeneralUnifier) FindSolution(t *term.Term) error {
	idx := mapToInt[t]
	rep := unif.ds.FindSet(idx)
	s := unif.schema[rep]
	if acyclic[rep] {
		return nil
	}
	if visited[rep] {
		return ErrUnifier
	}
	if s.Typ == term.TermCompound {
		visited[rep] = true
		for _, arg := range s.Args {
			unif.Initializer(arg, nil)
			if err := unif.FindSolution(arg); err != nil {
				return ErrUnifier
			}
		}
		visited[rep] = false
	}
	acyclic[rep] = true
	// Ensure that s has been registered.
	unif.Initializer(s, nil)
	rep2 := unif.ds.FindSet(mapToInt[s])
	varsList := unif.vars[rep2]
	for _, v := range varsList {
		if v != s {
			unifyMap[v] = s
		}
	}
	return nil
}
