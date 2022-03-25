package parser

import (
	"fmt"
	"github.com/butzopower/natsu/core"
	"github.com/butzopower/natsu/util"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
)

func Parse(unionPackageName, unionTypeName string) (core.Result, error) {
	var result core.Result

	obj, err := findType(unionPackageName, unionTypeName)

	if err != nil {
		return result, err
	}

	unionType, err := extractUnion(obj)

	if err != nil {
		return result, err
	}

	var terms []*types.Term

	for i := 0; i < unionType.Len(); i++ {
		terms = append(terms, unionType.Term(i))
	}

	termNames, err := util.MapWithErr(terms, func(term *types.Term) (core.TermPath, error) {
		full := term.String()
		pkg, local, splitErr := splitSourceType(full)

		if splitErr != nil {
			return core.TermPath{}, splitErr
		}

		return core.TermPath{
			Package: pkg,
			Local:   local,
		}, nil
	})

	if err != nil {
		return result, err
	}

	return core.Result{
		Path: unionPackageName,
		Union: core.TermPath{
			Package: unionPackageName,
			Local:   unionTypeName,
		},
		Terms: termNames,
	}, nil
}

func findType(sourceTypePackage, sourceTypeName string) (types.Object, error) {
	pkg, err := loadPackage(sourceTypePackage)

	if err != nil {
		return nil, err
	}

	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		return nil, fmt.Errorf("%s not found in declared types of %s", sourceTypeName, pkg)
	}
	return obj, nil
}

func extractUnion(obj types.Object) (*types.Union, error) {
	if _, ok := obj.(*types.TypeName); !ok {
		return nil, fmt.Errorf("%v is not a named type", obj)
	}

	interfaceType, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("type %v is not an interface", obj)
	}

	if interfaceType.NumEmbeddeds() == 0 {
		return nil, fmt.Errorf("type %v does not contain embedded types", obj)
	}

	unionType, ok := interfaceType.EmbeddedType(0).(*types.Union)
	if !ok {
		return nil, fmt.Errorf("type %v is not a union", obj)
	}

	return unionType, nil
}

func loadPackage(path string) (*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("loading packages for inspection: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0], nil
}

func splitSourceType(sourceType string) (string, string, error) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		return "", "", fmt.Errorf(`expected qualified type as "pkg/path.MyType"`)
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName, nil
}
