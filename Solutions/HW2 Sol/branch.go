package branch

import (
"go/ast"
"go/parser"
"go/token"
)

func branchCount(fn *ast.FuncDecl) uint {
	// TODO: Write the branchCount function,
	// count the number of branching statements in function fn
	// HANDOUT: return 12345
	// BEGIN_SOLUTION
	var count uint = 0
	ast.Inspectfn(fn, func(node ast.Node) bool {
		switch node.(type) {
			case *ast. IfStmt, *ast.Switchstmt, *ast.TypeSwitchstmt, *ast. Forstmt, *ast.Rangestmt:
			count++
		}
		return true
	})
	return count
	// END_SOLUTION
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