package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func ModelsData(f *ast.File, fset *token.FileSet) []Data {
	var results []Data
	ast.Inspect(f, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.TypeSpec:
			if t.Name == nil {
				return true
			}

			for _, c := range f.Comments {
				isValidComment := c.Pos() > t.Pos() && c.End() < t.End()
				if !isValidComment {
					continue
				}
				for _, l := range c.List {
					if strings.Contains(l.Text, magicString) {
						results = append(results, data(f, fset, t, l))
						return false
					}

				}
			}
		}
		return true
	})

	return results
}

func data(f *ast.File, fset *token.FileSet, t *ast.TypeSpec, comment *ast.Comment) Data {
	ff := fset.File(t.Pos())
	file, err := os.Open(ff.Name())
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}
	file.Close()
	src := buf.Bytes()
	fileData := Data{
		Type:    t.Name.Name,
		Package: f.Name.Name,
	}

	annotateFlagInfo(comment.Text, &fileData)
	ast.Inspect(t, func(node ast.Node) bool {
		switch t := node.(type) {
		case *ast.StructType:
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
				field.ParamName = LcFirst(field.FieldName) + "Param" // add suffix to avoid keyword collisions

				var isArray bool
				var isPointer bool

				start := fset.PositionFor(f.Type.Pos(), false)
				end := fset.PositionFor(f.Type.End(), false)
				rawType := string(src[start.Offset:end.Offset])
				if strings.HasPrefix(rawType, "*") {
					field.Star = "*"
					field.Point = "&"
					isPointer = true

					if strings.HasPrefix(rawType, "*[]") {
						isArray = true
						field.ParamType = "..." + rawType[3:]
					} else {
						field.ParamType = rawType[1:]
					}

				} else if strings.HasPrefix(rawType, "[]") {
					isArray = true
					isPointer = true
					field.ParamType = "..." + rawType[2:]
				}

				field.FieldType = rawType
				if field.ParamType == "" {
					field.ParamType = rawType
				}

				if isArray {
					field.FieldParamPrefix = "..."
				}

				fileData.BuilderFields = append(fileData.BuilderFields, field)
				if isPointer || isArray {
					fileData.PointerFields = append(fileData.PointerFields, field)
				} else {
					fileData.NonPointerFields = append(fileData.NonPointerFields, field)
				}

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

	for _, c := range f.Comments {
		for _, l := range c.List {
			if strings.Contains(l.Text, "+build ") {
				fileData.BuildTags += l.Text + "\n\n"
				break
			}

		}
		break
	}

	return fileData
}

type Data struct {
	BuildTags        string
	Package          string
	Imports          string
	Type             string
	BuilderFields    []BuilderField
	PointerFields    []BuilderField
	NonPointerFields []BuilderField
	Globals          bool
	NoBuilder        bool
	Prefix           string
	Suffix           string
}

func (d Data) FilePath(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s%s.go", filePrefix, ToSnakeCase(d.Type)))
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
