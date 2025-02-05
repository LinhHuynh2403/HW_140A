package branch

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func branchCount(fn *ast.FuncDecl) uint {
	var dem uint =0
	ast.Inspect(fn.Body, func(aa ast.Node) bool {
		switch aa.(type) {
		case *ast.IfStmt:
			dem++
		case *ast.SwitchStmt:
			dem++
		case *ast.TypeSwitchStmt:
			dem++
		case *ast.ForStmt:
			dem++
		case *ast.RangeStmt:
			dem++
		}
		return true
	})
	return dem
}

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
