package templates

import (
	"html/template"
	"os"
	"path/filepath"
)

func Load(root, ext string, funcMap template.FuncMap) (templates map[string]*template.Template, err error) {
	templates = make(map[string]*template.Template)
	funcMap["extend"] = func(_ string) string { return "" }

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ext || info.Name()[0] == '_' {
			return nil
		}

		name := path[len(root)+1:]
		name = filepath.ToSlash(name)

		t := template.New(name)
		t.Funcs(funcMap)

		err = Parse(t, path)
		if err == nil {
			templates[name] = t
		}
		return err
	})
	return
}
