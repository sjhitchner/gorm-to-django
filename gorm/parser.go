package gorm

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Model struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	Tags []Tag
}

type Tag struct {
	Name  string
	Value string
}

func Parse(packageDir string) error {
	// Create the file set and parse the package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packageDir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Loop through packages (assuming there's only one in the directory)
	for _, pkg := range pkgs {
		// Print package name
		fmt.Println("Package:", pkg.Name)

		// Loop through files in the package
		for filename, file := range pkg.Files {
			fmt.Println("\nFile:", filename)

			// Inspect the AST (Abstract Syntax Tree) of the file
			ast.Inspect(file, gormInspectV2(fset))
			// if typeSpec, ok := node.(*ast.TypeSpec); ok {
			//	fmt.Println("Type:", typeSpec.Name)
			//}

			/*
				// Print function declarations
				if funcDecl, ok := node.(*ast.FuncDecl); ok {
					fmt.Println("Function:", funcDecl.Name)
				}

				// Print variable declarations
				if varSpec, ok := node.(*ast.ValueSpec); ok {
					for _, name := range varSpec.Names {
						fmt.Println("Variable:", name)
					}
				}
			*/

		}
	}

	return nil
}

func gormInspectV2(fset *token.FileSet) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				fmt.Println("Struct:", typeSpec.Name)

				// Optionally, you can print fields of the struct
				for _, field := range structType.Fields.List {
					fieldName := field.Names[0].Name
					var fieldType string

					// Determine field data type
					switch t := field.Type.(type) {
					case *ast.Ident:
						fieldType = t.Name
					case *ast.SelectorExpr:
						fieldType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
					case *ast.StarExpr:
						if ident, ok := t.X.(*ast.Ident); ok {
							fieldType = "*" + ident.Name
						}
					case *ast.ArrayType:
						fieldType = fmt.Sprintf("[]%s", fset.Position(t.Pos()).String())
					case *ast.MapType:
						fieldType = fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())
					default:
						fieldType = "unknown"
					}

					var fieldTag string

					// Extract struct tags if present
					if field.Tag != nil {
						fieldTag = field.Tag.Value
					}

					fmt.Printf("  Field: %s %s %s\n", fieldName, fieldType, fieldTag)
				}
			}
		}

		return true
	}
}

func gormInspectV1(fset *token.FileSet) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		if typeSpec, ok := node.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				fmt.Println("Struct:", typeSpec.Name)

				// Optionally, you can print fields of the struct
				for _, field := range structType.Fields.List {
					fieldName := field.Names[0].Name
					var fieldType string

					// Determine field data type
					switch t := field.Type.(type) {
					case *ast.Ident:
						fieldType = t.Name
					case *ast.SelectorExpr:
						fieldType = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
					case *ast.StarExpr:
						if ident, ok := t.X.(*ast.Ident); ok {
							fieldType = "*" + ident.Name
						}
					case *ast.ArrayType:
						fieldType = fmt.Sprintf("[]%s", fset.Position(t.Pos()).String())
					case *ast.MapType:
						fieldType = fmt.Sprintf("map[%s]%s", fset.Position(t.Key.Pos()).String(), fset.Position(t.Value.Pos()).String())
					default:
						fieldType = "unknown"
					}

					var fieldTag string

					// Extract struct tags if present
					if field.Tag != nil {
						fieldTag = field.Tag.Value
					}

					fmt.Printf("  Field: %s %s %s\n", fieldName, fieldType, fieldTag)
				}
			}
		}

		return true
	}
}
