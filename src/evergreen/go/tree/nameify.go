package tree

import (
	"fmt"
	"strings"
)

type FileInfo struct {
	Package     *Package
	PackageName map[*Package]string
}

func DefaultPackageName(pkg *Package) string {
	n := len(pkg.Path)
	return pkg.Path[n-1]
}

func IsBuiltinPackage(pkg *Package) bool {
	return len(pkg.Path) == 0
}

func (info *FileInfo) ImportedName(pkg *Package) string {
	name, ok := info.PackageName[pkg]
	if !ok {
		// HACK assume no import name conflicts
		name = DefaultPackageName(pkg)
		info.PackageName[pkg] = name
	}
	return name
}

func (info *FileInfo) QualifyName(pkg *Package, name string) string {
	if pkg != info.Package && !IsBuiltinPackage(pkg) {
		return fmt.Sprintf("%s.%s", info.ImportedName(pkg), name)
	}
	return name
}

func nameifyType(t Type, info *FileInfo) {
	switch t := t.(type) {
	case *TypeRef:
		impl := t.Impl
		if impl == nil {
			panic(t)
		}
		switch impl := impl.(type) {
		case *StructDecl:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *InterfaceDecl:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *ExternalType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
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
	case *GetName, *GetLocal, *GetGlobal:
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

func nameifyTarget(expr Target, info *FileInfo) {
	switch expr := expr.(type) {
	case *SetName, *SetLocal:
		// TODO
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
			nameifyTarget(e, info)
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
		CompactFunc(decl)
		for _, lcl := range decl.Locals {
			nameifyType(lcl.T, info)
		}

		if decl.Recv != nil {
			nameifyType(decl.Recv.Type, info)
		}
		nameifyType(decl.Type, info)
		nameifyBody(decl.Body, info)
		InsertVarDecls(decl)
	default:
		panic(decl)
	}
}

func nameifyFile(pkg *Package, file *File) {
	file.Package = pkg.Path[len(pkg.Path)-1]

	info := &FileInfo{Package: pkg, PackageName: map[*Package]string{}}

	for _, decl := range file.Decls {
		nameifyDecl(decl, info)
	}

	// TODO clear existing imports.
	for pkg, name := range info.PackageName {
		path := strings.Join(pkg.Path, "/")
		file.Imports = append(file.Imports, &Import{Path: path, Name: name})
	}
}

func Nameify(prog *Program) {
	nameifyPrepass(prog)
	for _, pkg := range prog.Packages {
		if pkg.Extern {
			continue
		}
		pkgName := ""
		if len(pkg.Path) > 0 {
			pkgName = pkg.Path[len(pkg.Path)-1]
		}
		for _, file := range pkg.Files {
			file.Package = pkgName
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
				case *ExternalType:
					decl.Package = pkg
				default:
					panic(decl)
				}
			}
		}
	}
}
