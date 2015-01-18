package tree

import (
	"evergreen/compiler"
	"evergreen/go/core"
	"fmt"
	"strings"
)

type FileInfo struct {
	Package     *core.Package
	PackageName map[*core.Package]string
	Decl        *FuncDecl // HACK
	CoreProg    *core.CoreProgram
}

func DefaultPackageName(pkg *core.Package) string {
	n := len(pkg.Path)
	return pkg.Path[n-1]
}

func (info *FileInfo) ImportedName(pkg *core.Package) string {
	name, ok := info.PackageName[pkg]
	if !ok {
		// HACK assume no import name conflicts
		name = DefaultPackageName(pkg)
		info.PackageName[pkg] = name
	}
	return name
}

func (info *FileInfo) QualifyName(pkg *core.Package, name string) string {
	if pkg != nil && pkg != info.Package {
		return fmt.Sprintf("%s.%s", info.ImportedName(pkg), name)
	}
	return name
}

func (info *FileInfo) GetCallable(f core.Callable) Expr {
	switch f := f.(type) {
	case *core.Function:
		return &GetGlobal{
			Text: info.QualifyName(f.Package, f.Name),
		}
	case *core.IntrinsicFunction:
		return &GetGlobal{
			Text: f.Name,
		}
	default:
		panic(f)
	}
}

func (info *FileInfo) LocalName(lcl *LocalInfo) string {
	return lcl.Name
}

func nameifyParam(p *Param, info *FileInfo) {
	p.Name = info.LocalName(p.Info)
	p.Type = p.Info.T
	nameifyType(p.Type, info)
}

func nameifyType(t TypeRef, info *FileInfo) {
	switch t := t.(type) {
	case *NameRef:
		impl := t.T
		switch impl := impl.(type) {
		case *core.StructType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.InterfaceType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.ExternalType:
			t.Name = info.QualifyName(impl.Package, impl.Name)
		case *core.TypeDefType:
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
		f := expr.F
		if f != nil {
			expr.Expr = info.GetCallable(f)
		}
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
	case *IntLiteral, *Float32Literal, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral:
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

func isGoKeyword(s string) bool {
	switch s {
	case "break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "got", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var":
		return true
	default:
		return false
	}
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
		if count[name] > 1 || isGoKeyword(name) {
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

func nameifyFile(pkg *PackageAST, file *FileAST, coreProg *core.CoreProgram) {
	p := pkg.P
	file.Package = p.Path[len(p.Path)-1]

	info := &FileInfo{
		Package:     p,
		PackageName: map[*core.Package]string{},
		CoreProg:    coreProg,
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

func Nameify(status compiler.PassStatus, prog *ProgramAST, coreProg *core.CoreProgram) {
	status.Begin()
	defer status.End()

	nameifyPrepass(prog)
	for _, pkg := range prog.Packages {
		p := pkg.P
		pkgName := ""
		if len(p.Path) > 0 {
			pkgName = p.Path[len(p.Path)-1]
		}
		for _, file := range pkg.Files {
			file.Package = pkgName
			nameifyFile(pkg, file, coreProg)
		}
	}
}

func nameifyPrepass(prog *ProgramAST) {
	for _, pkg := range prog.Packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *InterfaceDecl:
					decl.T.Package = pkg.P
				case *StructDecl:
					decl.T.Package = pkg.P
				case *FuncDecl:
					decl.Package = pkg.P
				case *OpaqueDecl:
					decl.T.Package = pkg.P
				case *TypeDefDecl:
					decl.T.Package = pkg.P
				case *VarDecl:
					// Leaf
				default:
					panic(decl)
				}
			}
		}
	}
}
