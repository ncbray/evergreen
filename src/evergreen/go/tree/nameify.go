package tree

import (
	"fmt"
	"strings"
)

type FileInfo struct {
	Package     *PackageAST
	PackageName map[*PackageAST]string
	Decl        *FuncDecl // HACK
}

func DefaultPackageName(pkg *PackageAST) string {
	n := len(pkg.Path)
	return pkg.Path[n-1]
}

func IsBuiltinPackage(pkg *PackageAST) bool {
	// HACK pkg should not be nil
	return pkg == nil || len(pkg.Path) == 0
}

func (info *FileInfo) ImportedName(pkg *PackageAST) string {
	name, ok := info.PackageName[pkg]
	if !ok {
		// HACK assume no import name conflicts
		name = DefaultPackageName(pkg)
		info.PackageName[pkg] = name
	}
	return name
}

func (info *FileInfo) QualifyName(pkg *PackageAST, name string) string {
	if pkg != info.Package && !IsBuiltinPackage(pkg) {
		return fmt.Sprintf("%s.%s", info.ImportedName(pkg), name)
	}
	return name
}

func (info *FileInfo) LocalName(index LocalInfo_Ref) string {
	return info.Decl.LocalInfo_Scope.Get(index).Name
}

func nameifyParam(p *Param, info *FileInfo) {
	p.Name = info.LocalName(p.Info)
	nameifyType(p.Type, info)
}

func nameifyType(t TypeRef, info *FileInfo) {
	switch t := t.(type) {
	case *NameRef:
		impl := t.T
		switch impl := impl.(type) {
		case *StructType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *InterfaceType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *ExternalType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *TypeDefType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		default:
			panic(impl)
		}
	case *SliceRef:
		nameifyType(t.Element, info)
	case *PointerRef:
		nameifyType(t.Element, info)
	case *FuncTypeRef:
		for _, p := range t.Params {
			nameifyParam(p, info)
		}
		for _, r := range t.Results {
			nameifyParam(r, info)
		}
	default:
		panic(t)
	}
}

func nameifyExpr(expr Expr, info *FileInfo) {
	switch expr := expr.(type) {
	case *GetLocal, *GetName, *GetGlobal:
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
		stmt.Name = info.LocalName(stmt.Info)
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

func baseName(lcl *LocalInfo) string {
	letters := []rune(lcl.Name)
	end := len(letters)
	for end > 0 && letters[end-1] >= '0' && letters[end-1] <= '9' {
		end -= 1
	}
	name := string(letters[:end])
	if name == "" {
		name = "r"
	}
	return name
}

func nameifyFunc(decl *FuncDecl, info *FileInfo) {
	info.Decl = decl

	CompactFunc(decl)

	num_vars := decl.LocalInfo_Scope.Len()

	// Normalize variable names and count identical names.
	count := make(map[string]int, num_vars)
	iter := decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		name := baseName(lcl)
		lcl.Name = name
		count[name] += 1
	}

	// Disambiguate identical names.
	uids := make(map[string]int, num_vars)
	iter = decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		name := baseName(lcl)
		if count[name] > 1 {
			uid, _ := uids[name]
			uids[name] += 1
			lcl.Name = fmt.Sprintf("%s%d", name, uid)
		}
	}

	iter = decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		lcl := iter.Value()
		nameifyType(lcl.T, info)
	}
	if decl.Recv != nil {
		nameifyParam(decl.Recv, info)
	}
	nameifyType(decl.Type, info)
	nameifyBody(decl.Body, info)
	InsertVarDecls(decl)
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
	case *TypeDefDecl:
		nameifyType(decl.Type, info)
	case *FuncDecl:
		nameifyFunc(decl, info)
	case *VarDecl:
		nameifyType(decl.Type, info)
		nameifyExpr(decl.Expr, info)
	default:
		panic(decl)
	}
}

func nameifyFile(pkg *PackageAST, file *FileAST) {
	file.Package = pkg.Path[len(pkg.Path)-1]

	info := &FileInfo{
		Package:     pkg,
		PackageName: map[*PackageAST]string{},
	}

	for _, decl := range file.Decls {
		nameifyDecl(decl, info)
	}

	if len(file.Imports) != 0 {
		panic(file.Imports)
	}
	for pkg, name := range info.PackageName {
		path := strings.Join(pkg.Path, "/")
		file.Imports = append(file.Imports, &Import{Path: path, Name: name})
	}
}

func Nameify(prog *ProgramAST) {
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

func nameifyPrepass(prog *ProgramAST) {
	for _, pkg := range prog.Packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *InterfaceDecl:
					decl.T.Package = pkg
				case *StructDecl:
					decl.T.Package = pkg
				case *FuncDecl:
					decl.Package = pkg
				case *OpaqueDecl:
					decl.T.Package = pkg
				case *TypeDefDecl:
					decl.T.Package = pkg
				case *VarDecl:
					// Leaf
				default:
					panic(decl)
				}
			}
		}
	}
}
