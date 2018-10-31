package gowl

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template/parse"
)

func loadTemplates(root, ext string, funcMap template.FuncMap) (templates map[string]*template.Template, err error) {
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

		err = parseTemplate(t, path)
		if err == nil {
			templates[name] = t
		}
		return err
	})
	return
}

func parseTemplate(t *template.Template, path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = t.Parse(string(buf))
	if err != nil {
		return err
	}

	for _, node := range t.Tree.Root.Nodes {
		node, ok := node.(*parse.ActionNode)
		if !ok || len(node.Pipe.Cmds) != 1 {
			continue
		}
		args := node.Pipe.Cmds[0].Args
		if len(args) != 2 {
			continue
		}
		fn, ok := args[0].(*parse.IdentifierNode)
		if !ok || fn.Ident != "extend" {
			continue
		}
		extend, ok := args[1].(*parse.StringNode)
		if !ok || extend.Text == "" {
			return fmt.Errorf("unexpected parameter %s in extend clause", args[1])
		}
		p := extend.Text
		if strings.IndexByte("\\/", path[0]) == -1 {
			p = filepath.Join(filepath.Dir(path), p)
		}
		if err = parseTemplate(t, p); err != nil {
			return err
		}
	}
	return nil
}
