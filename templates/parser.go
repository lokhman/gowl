package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template/parse"
)

func Parse(t *template.Template, path string) error {
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
		if err = Parse(t, p); err != nil {
			return err
		}
	}
	return nil
}
