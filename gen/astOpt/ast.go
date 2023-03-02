package astOpt

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// AnalysisGo get ast from file
func AnalysisGo(Path string, Content string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	var err error
	var file *ast.File
	if Content == "" {
		file, err = parser.ParseFile(fset, Path, nil, parser.AllErrors)
	} else {
		file, err = parser.ParseFile(fset, Path, Content, parser.AllErrors)
	}
	if err != nil {
		return nil, nil, err
	}
	return fset, file, nil
}

// AstToGo write node to file
func AstToGo(f *os.File, node interface{}) error {
	var dst = bytes.Buffer{}

	err := astToBuf(&dst, node)
	if err != nil {
		return err
	}
	_, err = f.Write(dst.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func astToBuf(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

// Replace replace _TABLENAME_ to tablename
func Replace(funcDecl *ast.FuncDecl, from, to string) error {
	ast.Inspect(funcDecl, func(node ast.Node) bool {
		//fmt.Println(node.Pos())
		switch node := node.(type) {
		case *ast.Ident:
			identReplace(node, from, to)
		}
		return true
	})
	return nil
}

func identReplace(node *ast.Ident, from, to string) bool {
	node.Name = strings.Replace(node.Name, from, to, -1)
	return true
}
