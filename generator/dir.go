package generator

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
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
	filesToDelete := make(map[string]struct{})
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		include := strings.Index(info.Name(), filePrefix) != 0
		if !include {
			_, filename := path.Split(info.Name())
			filesToDelete[filename] = struct{}{}
		}

		return include
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
		output := m.FilePath(dir)
		_, filename := path.Split(output)
		fmt.Println("generating", filename, "for struct", m.Type)
		if err := DataToFile(output, tpl, m); err != nil {
			return err
		}

		delete(filesToDelete, filename)
	}

	for file := range filesToDelete {
		f := path.Join(dir, file)
		fmt.Println("deleting file", file, os.Remove(f))
	}

	return nil
}

func DataToFile(output string, tpl *template.Template, data Data) (err error) {
	defer func() {
		if err != nil {
			os.Remove(output)
		}
	}()

	//f, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	//if err != nil {
	//	return err
	//}
	//defer f.Close()

	buff := new(bytes.Buffer)

	err = tpl.Execute(buff, data)
	if err != nil {
		return err
	}
	//f.Seek(0, 0)
	//buf := new(bytes.Buffer)
	//io.Copy(buf, f)
	//f.Close()

	ff, err := imports.Process(output, buff.Bytes(), nil)
	if err != nil {
		return err
	}

	buff = bytes.NewBuffer(ff)

	f, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, buff)

	return err
}
