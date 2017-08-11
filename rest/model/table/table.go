package table

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

// ITable interface for implement table
type ITable interface {
	TableName() string
	CreateIfNotExist(db *sql.DB) (string, error)
	Query(db *sql.DB, condition string) ([]ITable, error)
	Insert(tx *sql.Tx) error
	Update(tx *sql.Tx) error
	Delete(tx *sql.Tx, condition string) error
}

var tables = make(map[string]ITable)

// Register register an implemented table
func Register(name string, table ITable) {
	if _, ok := tables[name]; !ok {
		tables[name] = table
	}
	//panic(fmt.Sprintf("table %s already registered", name))
}

//InitTables
func InitTables(db *sql.DB) {
	var sqls string
	for name, table := range tables {
		sql, err := table.CreateIfNotExist(db)
		if err != nil {
			panic(fmt.Errorf("create table %s --- %s", name, err.Error()))
		}
		sqls += fmt.Sprintln(sql)
	}
	ioutil.WriteFile("table.sql", []byte(sqls), 0666)
}

//DropTables
func DropTables(db *sql.DB) {
	sql := "DROP TABLE IF EXISTS `%s`;"
	for name := range tables {
		if _, err := db.Exec(fmt.Sprintf(sql, name)); err != nil {
			panic(err)
		}
	}
}

// Table
func Table(name string) (ITable, error) {
	if table, ok := tables[name]; ok {
		return table, nil
	}
	return nil, fmt.Errorf("not found table %s", name)
}
