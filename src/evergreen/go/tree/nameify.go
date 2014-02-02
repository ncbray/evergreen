package tree

import (
	"fmt"
)

type FileInfo struct {
	Package *Package
	Imports []string
}

func nameifyType(t Type, info *FileInfo) {
	switch t := t.(type) {
	case *TypeRef:
		impl := t.Impl
		if impl == nil {
			// TODO panic
			//fmt.Println("???", t.Name)
			return
		}
		switch impl := impl.(type) {
		case *StructDecl:
			if impl.Package != info.Package {
				panic(fmt.Sprintf("Can't handle interpackage refs: %#v vs %#v", impl.Package, info.Package))
			}
			t.Name = impl.Name
		case *InterfaceDecl:
			if impl.Package != info.Package {
				panic(fmt.Sprintf("Can't handle interpackage refs: %#v vs %#v", impl.Package, info.Package))
			}
			t.Name = impl.Name
		default:
			panic(impl)
		}
	case *SliceType:
		nameifyType(t.Element, info)
	case *PointerType:
		nameifyType(t.Element, info)
	case *FuncType:
		for _, p := range t.Params {
			nameifyType(p.Type, info)
		}
		for _, r := range t.Results {
			nameifyType(r.Type, info)
		}
	default:
		panic(t)
	}
}

func nameifyExpr(expr Expr, info *FileInfo) {
	switch expr := expr.(type) {
	case *NameRef:
		// TODO
	case *UnaryExpr:
		nameifyExpr(expr.Expr, info)
	case *BinaryExpr:
		nameifyExpr(expr.Left, info)
		nameifyExpr(expr.Right, info)
	case *Call:
		nameifyExpr(expr.Expr, info)
		for _, e := range expr.Args {
			nameifyExpr(e, info)
		}
	case *Selector:
		nameifyExpr(expr.Expr, info)
	case *Index:
		nameifyExpr(expr.Expr, info)
		nameifyExpr(expr.Index, info)
	case *TypeCoerce:
		nameifyType(expr.Type, info)
		nameifyExpr(expr.Expr, info)
	case *TypeAssert:
		nameifyExpr(expr.Expr, info)
		nameifyType(expr.Type, info)
	case *StructLiteral:
		nameifyType(expr.Type, info)
		for _, e := range expr.Args {
			nameifyExpr(e.Expr, info)
		}
	case *ListLiteral:
		nameifyType(expr.Type, info)
		for _, e := range expr.Args {
			nameifyExpr(e, info)
		}
	case *IntLiteral, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral:
		// Leaf
	default:
		panic(expr)
	}
}

func nameifyStmt(stmt Stmt, info *FileInfo) {
	switch stmt := stmt.(type) {
	case *Var:
		nameifyType(stmt.Type, info)
		if stmt.Expr != nil {
			nameifyExpr(stmt.Expr, info)
		}
	case *Assign:
		for _, e := range stmt.Sources {
			nameifyExpr(e, info)
		}
		for _, e := range stmt.Targets {
			nameifyExpr(e, info)
		}
	case *If:
		nameifyExpr(stmt.Cond, info)
		nameifyBody(stmt.Body, info)
		if stmt.Else != nil {
			nameifyStmt(stmt.Else, info)
		}
	case *BlockStmt:
		nameifyBody(stmt.Body, info)
	case *Return:
		for _, e := range stmt.Args {
			nameifyExpr(e, info)
		}
	case *Goto:
		// TODO
	case *Label:
		// TODO
	default:
		// TODO unhack.
		expr, ok := stmt.(Expr)
		if ok {
			nameifyExpr(expr, info)
		} else {
			panic(stmt)
		}
	}
}

func nameifyBody(body []Stmt, info *FileInfo) {
	for _, stmt := range body {
		nameifyStmt(stmt, info)
	}
}

func nameifyDecl(decl Decl, info *FileInfo) {
	switch decl := decl.(type) {
	case *InterfaceDecl:
		for _, field := range decl.Fields {
			nameifyType(field.Type, info)
		}
	case *StructDecl:
		for _, field := range decl.Fields {
			nameifyType(field.Type, info)
		}
	case *FuncDecl:
		if decl.Recv != nil {
			nameifyType(decl.Recv.Type, info)
		}
		nameifyType(decl.Type, info)
		nameifyBody(decl.Body, info)
	default:
		panic(decl)
	}
}

func nameifyFile(pkg *Package, file *File) {
	file.Package = pkg.Path[len(pkg.Path)-1]

	info := &FileInfo{Package: pkg}

	for _, decl := range file.Decls {
		nameifyDecl(decl, info)
	}
}

func Nameify(prog *Program) {
	nameifyPrepass(prog)
	for _, pkg := range prog.Packages {
		for _, file := range pkg.Files {
			nameifyFile(pkg, file)
		}
	}
}

func nameifyPrepass(prog *Program) {
	for _, pkg := range prog.Packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *InterfaceDecl:
					decl.Package = pkg
				case *StructDecl:
					decl.Package = pkg
				case *FuncDecl:
					decl.Package = pkg
				default:
					panic(decl)
				}
			}
		}
	}
}
