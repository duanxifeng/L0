package model

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/bocheninc/L0/rest/config"
	"github.com/bocheninc/L0/rest/model/database"
	//register database
	_ "github.com/bocheninc/L0/rest/model/database/mysql"
	"github.com/bocheninc/L0/rest/model/table"
	//register table
	_ "github.com/bocheninc/L0/rest/model/table/account"
	_ "github.com/bocheninc/L0/rest/model/table/history"
	_ "github.com/bocheninc/L0/rest/model/table/policy"
	_ "github.com/bocheninc/L0/rest/model/table/status"
	_ "github.com/bocheninc/L0/rest/model/table/transaction"
	_ "github.com/bocheninc/L0/rest/model/table/user"
)

var DB *sql.DB

func init() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8&loc=%s&parseTime=true",
		config.Cfg.DBUser, config.Cfg.DBPWD, config.Cfg.DBHost, config.Cfg.DBPort, url.QueryEscape(config.Cfg.DBZone))
	if tdb, err := database.Open(config.Cfg.DBEngine, connStr); err != nil {
		panic(fmt.Sprintf("open database error connStr=%s err=%v", connStr, err))
	} else {
		DB = tdb
	}

	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	if _, err := DB.Exec(fmt.Sprintf("create database if not exists %s;", config.Cfg.DBName)); err != nil {
		panic(err)
	}
	if _, err := DB.Exec(fmt.Sprintf("use %s;", config.Cfg.DBName)); err != nil {
		panic(err)
	}
	table.InitTables(DB)
}

func ClearDB() {
	table.DropTables(DB)
}

// //Execute
// func Execute(tableName, sql string) (table.ITable, error) {
// 	t, err := table.Table(tableName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return t.Execute(DB, sql)
// }

// //Update
// func Update(tableName, sql string) (table.ITable, error) {
// 	t, err := table.Table(tableName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return t.Update(DB, sql)
// }

// //Delete
// func Delete(tableName, sql string) (int, error) {
// 	t, err := table.Table(tableName)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return t.Delete(DB, sql)
// }

// //Query
// func Query(tableName, sql string) ([]table.ITable, error) {
// 	t, err := table.Table(tableName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return t.Query(DB, sql)
// }
