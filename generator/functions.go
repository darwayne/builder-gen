package generator

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
	"unicode"
)

func ModelsData(f *ast.File) []Data {
	var results []Data
	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.TypeSpec:
			if t.Name == nil {
				return true
			}

			for _, c := range f.Comments {
				for _, l := range c.List {
					if strings.Contains(l.Text, "::builder-gen") {
						results = append(results, data(f, t))
						return false
					}

				}
			}
		}
		return true
	})

	return results
}

func data(f *ast.File, t *ast.TypeSpec) Data {
	fileData := Data{
		Type:    t.Name.Name,
		Package: f.Name.Name,
	}
	ast.Inspect(t, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.StructType:
			_ = t
			fmt.Println("struct spotted")
			if t.Fields == nil {
				return false
			}

			for _, f := range t.Fields.List {
				if f == nil || len(f.Names) != 1 {
					continue
				}

				var field BuilderField
				field.FieldName = f.Names[0].Name
				field.FuncName = field.FieldName
				field.ParamName = LcFirst(field.FieldName)

				fmt.Println("field", f.Names[0].Name)
				//if f.Type != nil {
				//	ast.Inspect(f.Type, func(node ast.Node) bool {
				//		fmt.Printf("field:type::::::::%T\n", node)
				//		return true
				//	})
				//}
				var mainType string
				var isArray bool
				var isPointer bool
				ast.Inspect(f, func(node ast.Node) bool {
					switch t := node.(type) {
					case *ast.StarExpr:
						isPointer = true
						field.FieldType += "*"
					case *ast.Ident:
						if t.Obj != nil && t.Obj.Type != "var" {
							return false
						}
						field.FieldType += t.Name
						mainType += t.Name
						fmt.Println("____===", t.Name, t.Obj)
						if t.Obj != nil {
							fmt.Println("\t", t.Obj.Name)
						}
					case *ast.ArrayType:
						isArray = true
						field.FieldType += "[]"
						field.FieldParamPrefix = "..."
					}
					fmt.Printf(":::::::%T\n", node)
					return true
				})

				if isPointer {
					field.Star = "*"
					field.Point = "&"
				}
				if isArray {
					field.ParamType += "..."
				}
				field.ParamType += mainType
				fileData.BuilderFields = append(fileData.BuilderFields, field)

			}
			return false
		}
		return true
	})

	if len(f.Imports) > 0 {
		fileData.Imports = "import (\n"

		for _, s := range f.Imports {
			fileData.Imports += "\t"
			if s.Name != nil && s.Name.Name != "" {
				fileData.Imports += s.Name.Name + " "
			}
			fileData.Imports += s.Path.Value + "\n"
		}
		fileData.Imports += "\n)"
	}

	fileData.Package = f.Name.Name

	return fileData
}

type Data struct {
	Package       string
	Imports       string
	Type          string
	BuilderFields []BuilderField
}
type BuilderField struct {
	ParamName        string
	ParamType        string
	FuncName         string
	FieldName        string
	FieldType        string
	FieldParamPrefix string
	Star             string
	Point            string
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts the provided string to snake_case.
// Based on https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
func ToSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	return strings.ToLower(output)
}
