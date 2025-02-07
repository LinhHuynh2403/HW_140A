package term

import "fmt"

// TermType enumerates all types of terms
type TermType int

// Enumerates all types of terms.
const (
	TermAtom TermType = iota
	TermNumber
	TermVariable
	TermCompound
)

// Term type represents a term. As we will learn later in the course, such terms
// are used in Prolog.
type Term struct {
	// Typ is the type of this term
	Typ TermType

	// For atom, number and variable terms, the Literal field stores the literal
	// of this term, e.g. "123", "foo", "X". For compound terms, the Literal
	// field should be an empty string "".
	Literal string

	// For compound terms, the Functor field should be an atom term, and the
	// Args field should be a non-empty slice of its arguments. For non-compound
	// terms (viz., atoms, numbers and variables), the Functor field and the
	// Args field have to be nil.
	Functor *Term
	Args    []*Term
}

// String serializes a term into a string and implements the Stringer interface
// for the Term type.
// See also: https://tour.golang.org/methods/17
func (tm *Term) String() string {
	str := ""
	switch {
	case tm == nil:
	case tm.Typ == TermCompound:
		str = fmt.Sprintf("%s(%s)", tm.Functor, TermSliceToString(tm.Args))
	default:
		str = tm.Literal
	}
	return str
}

// TermSliceToString serializes a slice of terms (usually the Args field of a
// compound term) into a string.
func TermSliceToString(termSlice []*Term) string {
	str := ""
	for _, term := range termSlice {
		if str == "" {
			str = term.String()
		} else {
			str = fmt.Sprintf("%s, %s", str, term)
		}
	}
	return str
}
