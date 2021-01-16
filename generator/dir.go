package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Dir(dir string) error {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return strings.Index(info.Name(), "builder_gen_") != 0
	}, parser.ParseComments)
	if err != nil {
		return err
	}

	var models []Data
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			models = append(models, ModelsData(f)...)
		}
	}

	tpl, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	for _, m := range models {
		if err := DataToFile(dir, tpl, m); err != nil {
			return err
		}
	}

	fmt.Printf("%+v\n", models)

	return nil
}

func DataToFile(dir string, tpl *template.Template, data Data) (err error) {
	output := filepath.Join(dir, fmt.Sprintf("builder_gen_%s.go", ToSnakeCase(data.Type)))
	defer func() {
		if err != nil {
			os.Remove(output)
		}
	}()

	f, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tpl.Execute(f, data)

	return err
}
