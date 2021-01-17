package generator

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

type DirOpts struct {
	//::builder-gen
	Recursive           bool
	RecursiveExclusions []string
	Trace               bool
}

func Dir(dir string, opts ...DirOptsFunc) error {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	info := ToDirOpts(opts...)
	if info.Recursive {
		dirs := make(map[string]struct{})
		for _, ex := range info.RecursiveExclusions {
			f := strings.TrimSpace(ex)
			p, err := filepath.Abs(f)
			if err != nil {
				return err
			}

			dirs[p] = struct{}{}
		}
		return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info != nil && info.IsDir() {
				if strings.HasPrefix(info.Name(), ".") {
					dirs[path] = struct{}{}
					return nil
				}

				for ex := range dirs {
					if strings.HasPrefix(path, ex) {
						return nil
					}
				}

				return Dir(path)
			}

			return nil
		})
	}

	var mode = parser.ParseComments
	if info.Trace {
		mode |= parser.Trace
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
	}, mode)
	if err != nil {
		return err
	}

	var models []Data
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			models = append(models, ModelsData(f, fset)...)
		}
	}

	tpl, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}
	for _, m := range models {
		output := m.FilePath(dir)
		d, filename := path.Split(output)
		fmt.Println("generating", filename, "for struct", m.Type, "in", d)
		if err := DataToFile(output, tpl, m, info.Trace); err != nil {
			return err
		}

		delete(filesToDelete, filename)
	}

	for file := range filesToDelete {
		f := path.Join(dir, file)
		fmt.Println("deleting file", f, os.Remove(f))
	}

	return nil
}

func DataToFile(output string, tpl *template.Template, data Data, trace bool) (err error) {
	buff := new(bytes.Buffer)
	err = tpl.Execute(buff, data)
	if err != nil {
		return err
	}

	ff, err := imports.Process(output, buff.Bytes(), nil)
	if err != nil {
		if trace {
			fmt.Println("===", string(buff.Bytes()), "\n===")
		}
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
