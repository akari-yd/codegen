package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	//fset := token.NewFileSet()
	//path, err := filepath.Abs("hello.go")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = format.Node(os.Stdout, fset, file)
	//if err != nil {
	//	fmt.Println(err)
	//}
	ret := checkFunc("biz", "Biz.*", "mock", "Test.*")
	if ret {
		os.Exit(1)
	}
}

func checkFunc(sourceDir, sourceFun, targetDir, targetFun string) bool {
	functionSource := map[string]interface{}{}
	functionTarget := map[string]interface{}{}

	sourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		fmt.Println(err)
		return false
	}
	targetDir, err = filepath.Abs(targetDir)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = GetFunc(functionSource, sourceDir, sourceFun, true)
	if err != nil {
		fmt.Println(err)
	}
	err = GetFunc(functionTarget, targetDir, targetFun, true)
	if err != nil {
		fmt.Println(err)
	}
	mark := false
	for key := range functionSource {
		val, ok := functionTarget["Test"+key]
		if !ok || val == false {
			fmt.Printf("\033[1;31;40mfunction %s has no test function\033[0m\n", key)
			mark = true
		}
	}
	if !mark {
		fmt.Printf("\033[1;32;40mall function test\033[0m\n")
	}
	return mark
}

func GetFunc(functions map[string]interface{}, Dir string, reg string, op interface{}) error {
	fset := token.NewFileSet()
	Path, err := filepath.Abs(Dir)
	if err != nil {
		return nil
	}
	file, err := parser.ParseDir(fset, Path, nil, parser.AllErrors)
	if err != nil {
		return err
	}
	for _, f := range file {
		for _, file := range f.Files {
			for _, decl := range file.Decls {
				fun, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				ok, err := regexp.Match(reg, []byte(fun.Name.Name))
				if err != nil {
					return err
				}
				if ok {
					functions[fun.Name.Name] = op
				}
			}
		}
	}
	return nil
}
