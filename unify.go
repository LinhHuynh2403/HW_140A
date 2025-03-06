package unify

import (
	"errors"
 	"hw4/disjointset"
	"hw4/term"
)

// ErrUnifier is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrUnifier = errors.New("unifier error")

// UnifyResult is the result of unification. For example, for a variable term
// `s`, `UnifyResult[s]` is the term which `s` is unified with.
type UnifyResult map[*term.Term]*term.Term

// Unifier is the interface for the term unifier.
// Do not change the definition of this interface
type Unifier interface {
	Unify(*term.Term, *term.Term) (UnifyResult, error)
}

type GeneralUnifier struct{
	disjointsets map[*term.Term]disjointset.DisjointSet		// disjoint sets for each representitive
	size map[*term.Term]int									// size for each disjoint set
	schema map[*term.Term]*term.Term						// map from disjoint set to a Term
	vars map[*term.Term][]*term.Term						// representitive -> list of TermVariable
}

/******************** Global Variable Definition ******* ***************/

// global variable to store the return unify result
var unifyMap = UnifyResult{}

// map for node <-> int in order to apply disjointset
var mapToInt = map[*term.Term]int{}
var mapToTerm = map[int]*term.Term{}
var nodeCounter = 0

// visited flag for all node
var visited = map[*term.Term]bool{}

// acyclic flag for all nodes
var acyclic = map[*term.Term]bool{}

/****************************** End *************************************/


// reset the global variable
func resetGlobal() {
	// unifyMap = UnifyResult{}
	mapToInt = map[*term.Term]int{}
	mapToTerm = map[int]*term.Term{}
    nodeCounter = 0
	visited = map[*term.Term]bool{}
	acyclic = map[*term.Term]bool{}
}

// NewUnifier creates a struct of a type that satisfies the Unifier interface.
func NewUnifier() Unifier {
	var unifyObj Unifier = GeneralUnifier{ make(map[*term.Term]disjointset.DisjointSet), 
										   make(map[*term.Term]int),
										   make(map[*term.Term]*term.Term), 
									 	   make(map[*term.Term][]*term.Term) }
	return unifyObj 
}

func (unif GeneralUnifier) Initializer(t1 *term.Term, t2 *term.Term) {
	// maps  *Term <--->  int 
	if _, ok := mapToInt[t1]; !ok {
		mapToInt[t1] = nodeCounter
		mapToTerm[nodeCounter] = t1
		nodeCounter++
	}

	// create the disjointset class for two terms
	if _, ok := unif.disjointsets[t1]; !ok {
		set1 := disjointset.NewDisjointSet()
		unif.disjointsets[t1] = set1
		unif.size[t1] = 1		// initialize the size for the class
		unif.schema[t1] = t1	// intialize the schema for the class
		// initialize the vars for the term
		if t1.Typ == term.TermVariable {
			unif.vars[t1] = append(unif.vars[t1], t1)
		}
	}

	if t2 == nil {
		return
	}

	if _, ok := mapToInt[t2]; !ok {
		mapToInt[t2] = nodeCounter
		mapToTerm[nodeCounter] = t2
		nodeCounter++
	}

	if _, ok := unif.disjointsets[t2]; !ok {
		set2 := disjointset.NewDisjointSet()
		unif.disjointsets[t2] = set2
		unif.size[t2] = 1
		unif.schema[t2] = t2

		if t2 != nil && t2.Typ == term.TermVariable {
			unif.vars[t2] = append(unif.vars[t2], t2)
		}
	}
}

// implements the Unify method with UnionSet struct
func (unif GeneralUnifier) Unify(t1 *term.Term, t2 *term.Term) (UnifyResult, error) {
	// initialization
	unif.Initializer(t1, t2)
	unifyMap = UnifyResult{}

	err_clousure := unif.UnifClousure(t1, t2)
	if err_clousure != nil {
		resetGlobal()
		return nil, ErrUnifier
	}

	err_findSol := unif.FindSolution(t1)
	if err_findSol != nil {
		resetGlobal()
		return nil, ErrUnifier
	}
	
	return unifyMap, nil
}

func (unif GeneralUnifier) UnifClousure(t1 *term.Term, t2 *term.Term) error {
	num1 := unif.disjointsets[t1].FindSet(mapToInt[t1])		// num1 is int
	num2 := unif.disjointsets[t2].FindSet(mapToInt[t2])		// num2 is int return from FindSet()
	s := mapToTerm[num1]
	t := mapToTerm[num2]
	
	if s != t {
		schema_s := unif.schema[s]
		schema_t := unif.schema[t]
		if schema_s.Typ == term.TermVariable || schema_t.Typ == term.TermVariable {
			// one of their schema is variable 
			unif.Union(s, t)
		} else {
			// both are non-variable term
			if schema_s.Typ == term.TermCompound && schema_t.Typ == term.TermCompound {
				// both are compound term, compare the functor
				if schema_s.Functor == schema_t.Functor && len(schema_s.Args) == len(schema_t.Args){
					unif.Union(s, t)
					// loop through the args in two compound terms
					for i := range schema_s.Args {
						// create the disjointSet for parameters
						unif.Initializer(schema_s.Args[i], schema_t.Args[i])
						
						// try matching the parameters of compounds
						err := unif.UnifClousure(schema_s.Args[i], schema_t.Args[i])
						if err != nil {
							return ErrUnifier
						}
					}
				} else {
					// not match with functor
					return ErrUnifier
				}
			} else {
				// unmatched term type or single term
				return ErrUnifier
			}
		}

	}
	return nil
}

func (unif GeneralUnifier) Union(t1 *term.Term, t2 *term.Term) {
	num1 := unif.disjointsets[t1].FindSet(mapToInt[t1])
	s := mapToTerm[num1]

	num2 := unif.disjointsets[t2].FindSet(mapToInt[t2])
	t := mapToTerm[num2]

	if unif.size[s] >= unif.size[t] {
		unif.size[s] += unif.size[t]
		unif.vars[s] = append(unif.vars[s], unif.vars[t]...)	// append the vars(t) to vars(s)
		if unif.schema[s].Typ == term.TermVariable {
			unif.schema[s] = unif.schema[t]
		} 
		unif.disjointsets[s].UnionSet(mapToInt[s], mapToInt[t])
		unif.disjointsets[s] = disjointset.MergeToUnionSets(unif.disjointsets[s], unif.disjointsets[t])
		
		// iterate the children in tree B and 
		// update the disjointset[Bi] = disjointset[s]
		for key, _ := range disjointset.GetParent(unif.disjointsets[t]) {
			ti := mapToTerm[key]
			unif.disjointsets[ti] = unif.disjointsets[s]				// update disjointes for ti
		}
	} else {
		unif.size[t] += unif.size[s]
		unif.vars[t] = append(unif.vars[t], unif.vars[s]...)
		if unif.schema[t].Typ == term.TermVariable {
			unif.schema[t] = unif.schema[s]
		} 
		unif.disjointsets[t].UnionSet(mapToInt[t], mapToInt[s])
		unif.disjointsets[t] = disjointset.MergeToUnionSets(unif.disjointsets[t], unif.disjointsets[s])		
		
		for key, _ := range disjointset.GetParent(unif.disjointsets[s]) {
			si := mapToTerm[key]
			unif.disjointsets[si] = unif.disjointsets[t]				// update disjointes for ti
		}
	}
}

func (unif GeneralUnifier) FindSolution(t *term.Term) error {
	// fmt.Println(" ************** Inside the FindSolution ******************")
	// fmt.Println(" *** Debug info: original_s =",t ,"- from line 249")
	num := unif.disjointsets[t].FindSet(mapToInt[t])
	s := mapToTerm[num]
	s = unif.schema[s]

	if _, ok := acyclic[s]; ok {
	
	}
	// fmt.Println(" *** Debug info: visited[s] = ",visited[s] , "from line 262")
	if val, ok := visited[s]; ok {
		if val == true {
			return ErrUnifier
		}
	}

	if s.Typ == term.TermCompound {
		visited[s] = true
		for i := range s.Args {
			unif.Initializer(s.Args[i], nil)
			err := unif.FindSolution(s.Args[i])
			if err != nil {
				return ErrUnifier
			}		
		}
		visited[s] = false
	}
	acyclic[s] = true
	
	num2 := unif.disjointsets[t].FindSet(mapToInt[s])
	s2 := mapToTerm[num2]
	varsList := unif.vars[s2]

	if len(varsList) > 0 {
		for _, x := range varsList {
			if x != s {
				unifyMap[x] = s
			}
		}
	}
	return nil
}