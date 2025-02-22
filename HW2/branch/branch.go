package branch

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// count the number of branching statement
// brancghin statement are constructs in the code whree the program can take different paths of execution such as if, switch, for, etc
func branchCount(fn *ast.FuncDecl) uint {
	// TODO: Write the branchCount function,
	// count the number of branching statements in function fn
	// the goal of this function is to count how many of these branching statements exist in a given function
	var count uint

	// using Inspect to traverse
	// the function counts only the branching statement (if, switch, type switch, for, range)
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n.(type) {
		// here are the branch statement
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt:
			count++
		}
		return true
	})
	return count
}

// ComputeBranchFactors returns a map from the name of the function in the given
// Go code to the number of branching statements it contains.
func ComputeBranchFactors(src string) map[string]uint {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	m := make(map[string]uint)
	for _, decl := range f.Decls {
		switch fn := decl.(type) {
		case *ast.FuncDecl:
			m[fn.Name.Name] = branchCount(fn)
		}
	}

	return m
}
