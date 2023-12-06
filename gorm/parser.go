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
		return recurseType(fset, t.Type)

	case *ast.Ident:
		return t.Name

	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)

	case *ast.StarExpr:
		return "*" + recurseType(fset, t.X)

	case *ast.ArrayType:
		return "[]" + recurseType(fset, t.Elt)

	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())

	default:
		panic("unknown token")
	}
}

func parseTags(tag *ast.BasicLit) map[string]Tag {
	if tag == nil {
		return nil
	}

	values, ok := reflect.StructTag(tag.Value).Lookup("gorm")
	if !ok {
		return nil
	}

	tagMap := make(map[string]Tag)
	for _, tag := range strings.Split(values, ";") {
		value := strings.SplitN(tag, ":", 2)

		if value[0] == "" {
			continue
		}

		t := Tag{
			Name: value[0],
		}
		if len(value) > 1 {
			t.Value = value[1]
		}
		tagMap[t.Name] = t
	}
	return tagMap
}

/*
func parseField(fset *token.FileSet, field *ast.Field) Field {
	var fieldType string

	// Determine field data type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name

		if t.Obj != nil {
			switch d := t.Obj.Decl.(type) {
			case *ast.TypeSpec:
				fieldType = fmt.Sprintf("%v", d.Type)
			}
		}

	case *ast.SelectorExpr:
		fieldType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)

	case *ast.StarExpr:
		switch d := t.X.(type) {
		case *ast.Ident:
			fieldType = "*" + d.Name

		case *ast.SelectorExpr:
			fieldType = fmt.Sprintf("*%s.%s", d.X.(*ast.Ident).Name, d.Sel.Name)
		}

	case *ast.ArrayType:
		switch e := t.Elt.(type) {
		case *ast.StarExpr:
			fieldType = fmt.Sprintf("[]%v", e)

		default:
			fmt.Println("QQQQQ", t.Elt, reflect.TypeOf(t.Elt))
		}
		fieldType = fmt.Sprintf("[]%v", t.Elt)

	case *ast.MapType:
		fieldType = fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())

	default:
		fieldType = "unknown"
	}

	return Field{
		Name: field.Names[0].Name,
		Type: fieldType,
		Tags: parseTags(field.Tag),
	}
}

func gormInspect2(ch chan<- Struct, fset *token.FileSet) func(node ast.Node) bool {
	return func(node ast.Node) bool {

		switch t := node.(type) {
		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {

				// fields := make([]Field, 0, 20)
				for _, f := range structType.Fields.List {
					field := parseField(fset, f)
					fmt.Println(field)
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
		}

		return true
	}
}
*/
