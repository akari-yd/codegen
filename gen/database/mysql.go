package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jimsmart/schema"
	"github.com/smallnest/gen/dbmeta"
	"strings"
)

func MysqlCnt(sqlType, sqlConnStr, sqlTable, sqlDatabase string) (map[string]dbmeta.DbTableMeta, error) {
	db, err := sql.Open(sqlType, sqlConnStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	var dbTables []string
	if sqlTable != "" && sqlTable != "all" {
		dbTables = strings.Split(sqlTable, ",")
		//fmt.Printf("showing meta for table(s): %s\n", sqlTable)
	} else {
		//fmt.Printf("showing meta for all tables\n")
		dbTables_, err := schema.TableNames(db)
		for _, table := range dbTables_ {
			if table[0] == sqlDatabase {
				dbTables = append(dbTables, table[1])
			}
		}
		if err != nil {
			fmt.Printf("Error in fetching tables information from %s information schema from %s\n", sqlType, sqlConnStr)
			return nil, err
		}
	}

	tableInfos := make(map[string]dbmeta.DbTableMeta)
	for _, tableName := range dbTables {
		//fmt.Printf("---------------------------\n")
		if strings.HasPrefix(tableName, "[") && strings.HasSuffix(tableName, "]") {
			tableName = tableName[1 : len(tableName)-1]
		}

		//fmt.Printf("[%s]\n", tableName)

		tableInfo, err := dbmeta.LoadMeta(sqlType, db, sqlDatabase, tableName)
		if err != nil {
			fmt.Printf("Error getting table info for %s error: %v\n\n\n\n", tableName, err)
			continue
		}
		tableInfos[tableName] = tableInfo

		//fmt.Printf("\n\nDDL\n%s\n\n\n", tableInfo.DDL())
		//
		//for _, col := range tableInfo.Columns() {
		//	fmt.Printf("%s\n", col.String())
		//
		//	colMapping, err := dbmeta.SQLTypeToMapping(strings.ToLower(col.DatabaseTypeName()))
		//	if err != nil { // unknown type
		//		fmt.Printf("unable to find mapping for db type: %s\n", col.DatabaseTypeName())
		//		continue
		//	}
		//	fmt.Printf("     %s\n", colMapping.String())
		//}
		//primaryCnt := dbmeta.PrimaryKeyCount(tableInfo)
		//fmt.Printf("primaryCnt: %d\n", primaryCnt)
		//
		//fmt.Printf("\n\n")
		//delSQL, err := dbmeta.GenerateDeleteSQL(tableInfo)
		//if err == nil {
		//	fmt.Printf("delSQL: %s\n", delSQL)
		//}
		//
		//updateSQL, err := dbmeta.GenerateUpdateSQL(tableInfo)
		//if err == nil {
		//	fmt.Printf("updateSQL: %s\n", updateSQL)
		//}
		//
		//insertSQL, err := dbmeta.GenerateInsertSQL(tableInfo)
		//if err == nil {
		//	fmt.Printf("insertSQL: %s\n", insertSQL)
		//}
		//
		//selectOneSQL, err := dbmeta.GenerateSelectOneSQL(tableInfo)
		//if err == nil {
		//	fmt.Printf("selectOneSQL: %s\n", selectOneSQL)
		//}
		//
		//selectMultiSQL, err := dbmeta.GenerateSelectMultiSQL(tableInfo)
		//if err == nil {
		//	fmt.Printf("selectMultiSQL: %s\n", selectMultiSQL)
		//}
	}

	return tableInfos, nil
}
