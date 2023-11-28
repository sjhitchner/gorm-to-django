package gorm

import (
	"encoding/json"
	"fmt"
	"go/ast"
	// "go/doc"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type Struct struct {
	IsModel  bool
	Name     string
	Metadata map[string]string
	Fields   []Field
}

func (t Struct) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(b)
}

type Field struct {
	Name string
	Type string
	Kind string
	Tags []Tag
}

type Tag struct {
	Name  string
	Value string
}

func Parse(packageDir string) (<-chan Struct, error) {

	out := make(chan Struct)

	// Create the file set and parse the package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packageDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)

		// Loop through packages (assuming there's only one in the directory)
		for _, pkg := range pkgs {
			// Print package name
			fmt.Println("Package:", pkg.Name)

			// Loop through files in the package
			for filename, file := range pkg.Files {
				fmt.Println("\nFile:", filename)

				// Inspect the AST (Abstract Syntax Tree) of the file
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

				fmt.Println("Struct:", t.Name)
				fields := make([]Field, 0, 20)

				// Optionally, you can print fields of the struct
				for _, f := range structType.Fields.List {
					field := parseField(fset, f)
					fmt.Printf(" Field: %v\n", field)
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

	// Split the input string by spaces to get key-value pairs
	keyValuePairs := strings.Fields(input)

	// Iterate over each key-value pair
	for _, pair := range keyValuePairs {
		// Split the pair by the colon to separate key and value
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
	var fieldType, fieldKind string

	// Determine field data type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fmt.Println("ident", t.Name)
		fieldType = t.Name
		fieldKind = t.Name

		if t.Obj != nil {
			switch d := t.Obj.Decl.(type) {
			case *ast.TypeSpec:
				fieldKind = fmt.Sprintf("%v", d.Type)
			}
		}

	case *ast.SelectorExpr:
		fmt.Println("times", t.X.(*ast.Ident).Name)
		fieldType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)

	case *ast.StarExpr:

		switch d := t.X.(type) {
		case *ast.Ident:
			fieldType = "*" + d.Name
			fieldKind = "struct"

		case *ast.SelectorExpr:
			fieldType = fmt.Sprintf("*%s.%s", d.X.(*ast.Ident).Name, d.Sel.Name)
			fieldKind = fmt.Sprintf("%s.%s", d.X.(*ast.Ident).Name, d.Sel.Name)
		}

	case *ast.ArrayType:
		// fieldType = fmt.Sprintf("[]%s", fset.Position(t.Pos()).String())
		fieldType = fmt.Sprintf("[]%v", t.Elt)
		fieldKind = "array"

	case *ast.MapType:
		fieldType = fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())
		fieldKind = "map"

	default:
		fieldType = "unknown"
	}

	return Field{
		Name: field.Names[0].Name,
		Type: fieldType,
		Kind: fieldKind,
		Tags: parseTags(field.Tag),
	}
}

func parseTags(tag *ast.BasicLit) []Tag {
	if tag == nil {
		return nil
	}

	values, ok := reflect.StructTag(tag.Value).Lookup("gorm")
	if !ok {
		return nil
	}

	tags := make([]Tag, 0, 5)
	for _, tag := range strings.Split(values, ";") {
		value := strings.SplitN(tag, ":", 2)
		if len(value) > 1 {
			tags = append(tags, Tag{
				Name:  value[0],
				Value: value[1],
			})
		} else {
			tags = append(tags, Tag{
				Name: value[0],
			})
		}
	}
	return tags
}
