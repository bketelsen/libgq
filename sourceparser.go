package libgq

import (
	"fmt"
	"go/importer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bketelsen/libgq/ast"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
)

type PackageParser struct {
	PkgName string
	Schema  string
}

func NewPackageParser(pkgName string) *PackageParser {
	return &PackageParser{
		PkgName: pkgName,
		Schema:  "",
	}
}

func (sp *PackageParser) Parse() error {
	defs := make(map[string]*ast.Definition)
	var defNames []string
	fset := token.NewFileSet()
	pkg, err := importer.ForCompiler(fset, "source", nil).Import(sp.PkgName)
	//pkg, err := importer.Default().Import(sp.PkgName)
	if err != nil {
		return err
	}
	for _, declName := range pkg.Scope().Names() {
		if shouldParse(declName) {
			fmt.Println(declName)
			defNames = append(defNames, declName)
			defs[declName] = &ast.Definition{Name: declName}
			/*obj := pkg.Scope().Lookup(declName)
			fmt.Println("exported: ", obj.Exported())
			fmt.Println("id: ", obj.Id())
			fmt.Println("name: ", obj.Name())
			fmt.Println("parent: ", obj.Parent())
			fmt.Println("package", obj.Pkg())
			fmt.Println("pos: ", obj.Pos())
			fmt.Println("string: ", obj.String())
			fmt.Println("type: ", obj.Type())
			*/
		}
	}
	gproot := os.Getenv("GOPATH")
	cfg := &packages.Config{
		//Dir:  "/home/bketelsen/projects/gq/src/blog/models",
		Dir:  gproot,
		Mode: packages.LoadAllSyntax}

	pkgs, err := packages.Load(cfg, sp.PkgName)
	if err != nil {
		return err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return err
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		fmt.Println("Package::: ", pkg.Name, pkg.ID, pkg.GoFiles)
		d := decorator.NewDecoratorFromPackage(pkg)
		for _, fname := range pkg.GoFiles {
			fmt.Println("\n\n\tFile: ", fname)
			src, err := ioutil.ReadFile(fname)

			if err != nil {
				return err
			}
			f, err := d.Parse(src)
			if err != nil {
				return err
			}
			model := filepath.Base(fname)
			model = strings.Replace(model, ".go", "", -1)
			model = strings.Title(model)
			fmt.Println(model)
			dst.Inspect(f, func(x dst.Node) bool {
				s, ok := x.(*dst.StructType)
				if !ok {
					return true
				}
				fl := []*ast.FieldDefinition{}

				for _, field := range s.Fields.List {
					fmt.Printf("Field: %s\n", field.Names[0].Name)
					fmt.Printf("Type: %v\n", field.Type)

					fmt.Printf("Type: %T\n", field.Type)
					fmt.Printf("Tag:   %s\n", field.Tag.Value)
					fd := &ast.FieldDefinition{
						Name: field.Names[0].Name,
						Type: ast.NamedType(fieldType(field), nil),
					}
					fmt.Println(fieldType(field))
					fl = append(fl, fd)
				}
				o := &ast.Definition{
					Kind:   ast.Object,
					Name:   model,
					Fields: fl,
				}

				defs[model] = o
				fmt.Println(defs[model].String())
				sp.Schema = sp.Schema + "\n" + defs[model].String() + "\n"
				return false
			})
		}

	}

	return nil

}

func fieldType(f *dst.Field) string {
	switch t := f.Type.(type) {
	case *dst.Ident:
		return strings.Title(t.String())

	case *dst.SelectorExpr:
		if t.Sel.String() == "UUID" {
			return "ID"
		}

		if t.Sel.String() == "Time" {
			return "Int"
		}
		return t.Sel.String()
	default:
		return ""
	}
}

// TODO: this is buffalo/pop specific,
// should be configurable
var skiplist = []string{"DB", "init"}

func shouldParse(name string) bool {
	for _, skip := range skiplist {
		if name == skip {
			return false
		}
	}
	return true
}
