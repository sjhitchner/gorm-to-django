package gorm

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

func Parse(packageDir string) (<-chan Struct, error) {

	if packageDir == "" {
		return nil, fmt.Errorf("No package directory given")
	}

	out := make(chan Struct)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packageDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				ast.Inspect(file, gormInspect(out, fset))
			}
		}
	}()

	return out, nil
}

func gormInspect(ch chan<- Struct, fset *token.FileSet) func(node ast.Node) bool {
	var comments []string

	return func(node ast.Node) bool {
		var comment string

		switch t := node.(type) {
		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {
				comment, comments = pop(comments)

				fields := make([]Field, 0, 20)
				for _, f := range structType.Fields.List {
					field := parseField(fset, f)
					fields = append(fields, field)
				}

				ch <- Struct{
					IsModel:  strings.HasPrefix(comment, "g2d"),
					Name:     t.Name.String(),
					Metadata: parseComment(comment),
					Fields:   fields,
				}
			}

		case *ast.GenDecl:
			comment = t.Doc.Text()
			if comment != "" {
				comments = append(comments, comment)
			}
		}
		return true
	}
}

func parseComment(input string) map[string]string {
	result := make(map[string]string)

	keyValuePairs := strings.Fields(input)
	for _, pair := range keyValuePairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			result[key] = value
		}
	}

	return result
}

func pop(s []string) (string, []string) {
	if len(s) > 0 {
		return s[len(s)-1], s[:len(s)-1]
	}
	return "", s
}

func parseField(fset *token.FileSet, field *ast.Field) Field {
	typ := recurseType(fset, field)

	return Field{
		Name: field.Names[0].Name,
		Type: typ,
		Tags: parseTags(field.Tag),
	}
}

func recurseType(fset *token.FileSet, node ast.Node) string {
	switch t := node.(type) {
	case *ast.Field:
		// fmt.Println("field", t)

		if t.Tag != nil {
			// fmt.Println("kind", t.Tag.Kind)
		}

		return recurseType(fset, t.Type)

	case *ast.Ident:
		// fmt.Println("ident", t.NamePos, t.Name)
		if t.Obj != nil {
			if n, ok := t.Obj.Decl.(ast.Node); ok {
				//fmt.Println("obj", t.Obj.Name, t.Obj.Decl, reflect.TypeOf(t.Obj.Decl))
				return recurseType(fset, n)
			}
		}
		return t.Name

	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)

	case *ast.StarExpr:
		return "*" + recurseType(fset, t.X)

	case *ast.ArrayType:
		return "[]" + recurseType(fset, t.Elt)

	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())

	case *ast.TypeSpec:
		switch t.Type.(type) {
		case *ast.StructType:
			return t.Name.Name
		default:
			return recurseType(fset, t.Type)
		}

	default:
		panic(fmt.Errorf("unknown token: %s", reflect.TypeOf(t)))
	}
}

func parseTags(tag *ast.BasicLit) map[string]Tag {
	if tag == nil {
		return nil
	}

	tagMap := make(map[string]Tag)
	parseTagByName(tag, tagMap, "gorm")
	parseTagByName(tag, tagMap, "django")
	return tagMap
}

func parseTagByName(tag *ast.BasicLit, tagMap map[string]Tag, name string) {
	values, ok := reflect.StructTag(tag.Value).Lookup(name)
	if !ok {
		return
	}

	for _, tag := range strings.Split(values, ";") {
		value := strings.SplitN(tag, ":", 2)

		if value[0] == "" {
			continue
		}

		t := Tag{
			Name:   value[0],
			Source: name,
		}
		if len(value) > 1 {
			t.Value = value[1]
		}
		tagMap[t.Name] = t
	}
}
