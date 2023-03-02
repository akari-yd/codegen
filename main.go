package main

import (
	"fmt"
	"git.garena.com/dang.yang/codegen/gen/astOpt"
	"git.garena.com/dang.yang/codegen/gen/database"
	"github.com/droundy/goopt"
	"github.com/go-toolsmith/astcopy"
	"github.com/go-toolsmith/astequal"
	"github.com/gobuffalo/packr/v2"
	"github.com/iancoleman/strcase"
	"github.com/smallnest/gen/dbmeta"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	sqlType       = goopt.String([]string{"--sqltype"}, "mysql", "sql database type such as [ mysql, mssql, postgres, sqlite, etc. ]")
	sqlConnStr    = goopt.String([]string{"-c", "--connstr"}, "", "database connection string")
	sqlDatabase   = goopt.String([]string{"-d", "--database"}, "", "Database to for connection")
	sqlTable      = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
	targetPath    = goopt.String([]string{"-T", "--target"}, "graph", "Generate path to project")
	sourcePath    = goopt.String([]string{"-s", "--source"}, "codegenTemplate", "Source path of template")
	force         = goopt.Flag([]string{"-f", "--force"}, nil, "Force to cover function with changed", "")
	filelist      = goopt.String([]string{"-F", "--files"}, "", "Select template to build,use \"all\" for all template")
	option        = goopt.String([]string{"-o", "--option"}, "biz,schema,tools", "Choose part of project to generate, include biz,schema,tools")
	baseTemplates *packr.Box
	dbtograph     map[string]string
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "ORM and RESTful meta data viewer for SQl databases"
	}
	goopt.Version = "v1.0.0 (2022/5/27)"
	goopt.Summary = `codegen [-v] --sqltype=mysql --connstr "user:password@/dbname" --database <databaseName> 

           sqltype - sql database type such as [ mysql, mssql, postgres, sqlite, etc. ]

`

	//Parse options
	goopt.Parse(nil)

}

func main() {
	//type of database to type of graphqls
	dbtograph = make(map[string]string)
	dbtograph["int"] = "Int"
	dbtograph["uint"] = "Int"
	dbtograph["bigint"] = "Int"
	dbtograph["ubigint"] = "Int"
	dbtograph["tinyint"] = "Int"
	dbtograph["utinyint"] = "Int"
	dbtograph["text"] = "String"
	dbtograph["varchar"] = "String"

	confList := []string{"makefile", "readme.md", "tools/makefile", "tools/readme.md", "tools/tools/.golangci.yml", "tools/tools/ast.go"}

	options := strings.Split(*option, ",")
	Project := false
	Schema := false
	Tools := false

	for _, op := range options {
		if op == "biz" {
			Project = true
		} else if op == "schema" {
			Schema = true
		} else if op == "tools" {
			Tools = true
		} else {
			fmt.Println("no option like ", op)
		}
	}

	baseTemplates = packr.New("template", *sourcePath)
	var tables map[string]dbmeta.DbTableMeta
	var err error
	if Project || Schema {
		tables, err = database.MysqlCnt(*sqlType, *sqlConnStr, *sqlTable, *sqlDatabase)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.Stat(*targetPath)
		if err != nil || !f.IsDir() {
			os.Mkdir(*targetPath, os.ModePerm)
		}

		if *sqlTable == "" {
			*sqlTable = "all"
		}
	}

	if Project {
		var goFiles []string
		if *filelist == "" || strings.EqualFold(*filelist, "all") {
			goFiles = []string{"biz/init.go", "biz/prj.db.go"}
		} else {
			goFiles = strings.Split(*filelist, ",")
		}

		for _, file := range goFiles {
			err = GenFile(targetPath, file, &tables)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	if Schema {
		err = GenGraphql(targetPath, "schema/prj.db.graphqls", tables)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if Tools {
		err = GenConf(confList, "")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// GenFile generate file from template,if tables is not nil,replace TABLENAME to table's real name
func GenFile(targetPath *string, filename string, tables *map[string]dbmeta.DbTableMeta) error {
	var temp *ast.File
	tempContent, err := baseTemplates.Find(filename)
	if err != nil {
		return err
	}
	_, temp_, err := astOpt.AnalysisGo("", string(tempContent))
	if err != nil {
		return err
	}
	if tables != nil && filename == "biz/prj.db.go" {
		temp = GenFromTable(temp_, *tables)
	} else {
		temp = temp_
	}

	file := *targetPath + "/" + filename
	dir := filepath.Dir(file)
	if f, err := os.Stat(dir); err != nil || !f.IsDir() {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	var gen *ast.File
	finfo, err := os.Stat(file)

	if err != nil || finfo.IsDir() {
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		gen = astcopy.File(temp)
		err = astOpt.AstToGo(f, gen)
		if err != nil {
			return err
		}
		f.Close()
	} else {
		fset, gen, err := astOpt.AnalysisGo(file, "")
		if err != nil {
			return err
		}
		//we should do follow clause after AnalysisGo
		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		absPath, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		FuncCmp(temp, gen, absPath, fset)
		err = astOpt.AstToGo(f, gen)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// FuncCmp compare all function, if function has been change, notice!
func FuncCmp(initTemp, initGen *ast.File, absPath string, fset *token.FileSet) {
	for _, declTemp := range initTemp.Decls {
		if funcDeclTemp, ok := declTemp.(*ast.FuncDecl); ok {
			equal := false //mark if find same name function
			for i, declGen := range initGen.Decls {
				if astequal.Decl(declTemp, declGen) {
					equal = true
					break
				} else if funcDeclGen, ok := declGen.(*ast.FuncDecl); ok {
					if astequal.Node(funcDeclGen.Name, funcDeclTemp.Name) {
						if *force {
							initGen.Decls[i] = funcDeclTemp
						} else {
							fmt.Printf("function %s:%d:function %s has been change\n", absPath, fset.Position(funcDeclGen.Pos()).Line, funcDeclGen.Name.Name)
						}
						equal = true
						break
					} else {
						continue
					}
				} else {
					continue
				}
			}
			if equal == false {
				initGen.Decls = append(initGen.Decls, astcopy.FuncDecl(funcDeclTemp))
			}
		}
	}
}

// GenFromTable generate ast from template and replace TABLENAME to table's real name
func GenFromTable(source *ast.File, tables map[string]dbmeta.DbTableMeta) (target *ast.File) {
	target = &ast.File{
		Doc:      source.Doc,
		Package:  source.Package,
		Name:     source.Name,
		Scope:    source.Scope,
		Comments: source.Comments,
	}
	for _, decl := range source.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			for tablename := range tables {
				name := strcase.ToCamel(tablename)
				fun := astcopy.FuncDecl(funcDecl)
				err := astOpt.Replace(fun, "_TABLENAME_", name)
				if err != nil {
					return nil
				}
				target.Decls = append(target.Decls, fun)
			}
		} else {
			target.Decls = append(target.Decls, decl)
		}
	}
	return target
}

func GenGraphql(targetPath *string, filename string, tables map[string]dbmeta.DbTableMeta) error {
	dir := filepath.Dir(*targetPath + "/" + filename)
	if f, err := os.Stat(dir); err != nil || !f.IsDir() {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(*targetPath+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	types := ""
	mutation := "type Mutation {\n"
	query := "type Query {\n"
	for tablename, meta := range tables {
		tablename = strcase.ToCamel(tablename)
		types += fmt.Sprintf("type %s {\n", tablename)
		for _, col := range meta.Columns() {
			types += fmt.Sprintf("\t\"\"\"%s\"\"\"\n", col.Comment())
			types += fmt.Sprintf("\t%s:%s\n\n", col.Name(), dbtograph[strings.ToLower(col.ColumnType())])
		}
		types += "}\n"

		types += fmt.Sprintf("input %sInput {\n", tablename)
		for _, col := range meta.Columns() {
			types += fmt.Sprintf("\t\"\"\"%s\"\"\"\n", col.Comment())
			types += fmt.Sprintf("\t%s:%s\n\n", col.Name(), dbtograph[strings.ToLower(col.ColumnType())])
		}
		types += "}\n\n"

		types += fmt.Sprintf("type %sPagerInfo {\n", tablename)
		types += fmt.Sprintf("\titems:[%s!]", tablename)
		types += "total:Int!\n}\n\n"

		mutation += fmt.Sprintf("\tcreate%s(input:%sInput!):%s!\n", tablename, tablename, tablename)
		mutation += fmt.Sprintf("\tupdate%s(id:Int!,input:%sInput!):%s!\n", tablename, tablename, tablename)
		mutation += fmt.Sprintf("\t#delete%s(id:Int!):Boolean!\n\n", tablename)

		query += fmt.Sprintf("\tget%ssForPager(query:%sInput,pager:PagerInput!,order:String):%sPagerInfo!\n", tablename, tablename, tablename)
		query += fmt.Sprintf("\tget%ss(query:%sInput,pager:PagerInput!,order:String):[%s]!\n", tablename, tablename, tablename)
		query += fmt.Sprintf("\tget%s(id:Int,query:%sInput):%s!\n\n", tablename, tablename, tablename)
	}

	mutation += "}\n\n"
	query += "}\n\n"

	types += "directive @goModel(\n    model: String\n    models: [String!]\n) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION\n\ndirective @goField(\n    forceResolver: Boolean\n    name: String\n) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION\n        \ndirective @goTag(\n    key: String!\n    value: String\n) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION\ninput PagerInput {\n    limit:Int!\n    offset:Int!\n}\n"

	_, err = file.WriteString(types)
	if err != nil {
		return err
	}
	_, err = file.WriteString(query)
	if err != nil {
		return err
	}
	_, err = file.WriteString(mutation)
	if err != nil {
		return err
	}

	return nil
}

func GenConf(confList []string, target string) error {
	stat, err := os.Stat(target + "tools")
	if err != nil || !stat.IsDir() {
		err := os.Mkdir(target+"/"+"tools", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		err = os.Mkdir(target+"/"+"tools/tools", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, conf := range confList {
		content, err := baseTemplates.Find(conf)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(target+"/"+conf, os.O_CREATE|os.O_WRONLY, 0666)
		defer file.Close()
		if err != nil {
			return err
		}
		_, err = file.Write(content)
		if err != nil {
			return err
		}
		if conf == "makefile" {
			_, err = file.WriteString(genCommand())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// generate command into make file todo make biz, make schema, make tools
func genCommand() (ret string) {
	regen := "regen: ## regen project\n"
	if *sqlType != "" {
		regen += "\t@codegen --sqltype=\"" + *sqlType + "\" \\\n"
	}
	if *sqlConnStr != "" {
		regen += "\t-c \"" + *sqlConnStr + "\" \\\n"
	}
	if *sqlDatabase != "" {
		regen += "\t-d \"" + *sqlDatabase + "\" \\\n"
	}
	if *sqlTable != "" {
		regen += "\t-t \"" + *sqlTable + "\" \\\n"
	}
	if *targetPath != "" {
		regen += "\t-T " + *targetPath + " \\\n"
	}
	if *sourcePath != "" {
		regen += "\t-s " + *sourcePath + " \\\n"
	}
	if *filelist != "" {
		regen += "\t-F \"" + *filelist + "\" \\\n"
	}
	if *force {
		regen += " \t-f \\\n"
	}
	regen += "\n"

	biz := "biz: ## generate biz project \n"
	if *sqlType != "" {
		biz += "\t@codegen --sqltype=\"" + *sqlType + "\" \\\n"
	}
	if *sqlConnStr != "" {
		biz += "\t-c \"" + *sqlConnStr + "\" \\\n"
	}
	if *sqlDatabase != "" {
		biz += "\t-d \"" + *sqlDatabase + "\" \\\n"
	}
	if *sqlTable != "" {
		biz += "\t-t \"" + *sqlTable + "\" \\\n"
	}
	if *targetPath != "" {
		biz += "\t-T " + *targetPath + " \\\n"
	}
	if *sourcePath != "" {
		biz += "\t-s " + *sourcePath + " \\\n"
	}
	if *filelist != "" {
		biz += "\t-F \"" + *filelist + "\" \\\n"
	}
	if *force {
		biz += " \t-f \\\n"
	}
	biz += "\n"

	Schema := "schema: ## make schemas \n"
	if *sqlType != "" {
		Schema += "\t@codegen --sqltype=\"" + *sqlType + "\" \\\n"
	}
	if *sqlConnStr != "" {
		Schema += "\t-c \"" + *sqlConnStr + "\" \\\n"
	}
	if *sqlDatabase != "" {
		Schema += "\t-d \"" + *sqlDatabase + "\" \\\n"
	}
	if *sqlTable != "" {
		Schema += "\t-t \"" + *sqlTable + "\" \\\n"
	}
	if *targetPath != "" {
		Schema += "\t-T " + *targetPath + " \\\n"
	}
	if *sourcePath != "" {
		Schema += "\t-s " + *sourcePath + " \\\n"
	}
	if *force {
		Schema += " \t-f \\\n"
	}
	Schema += "\n"

	Tools := "tools: ## generate tools like makefile tools/.\n"
	Tools += "\t@codegen -o \"tools\"\n\n"

	ret = "\n\n"

	if *option == "biz,schema,tools" {
		ret += regen + biz + Schema + Tools
	} else {
		ret += Tools
	}
	return ret
}
