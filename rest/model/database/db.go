package database

import (
	"database/sql"
	"fmt"
)

//IDataBase interface for implement database
type IDataBase interface {
	Open(string) (*sql.DB, error)
}

var databases = make(map[string]IDataBase)

//Register register an implemented database
func Register(name string, database IDataBase) {
	if _, ok := databases[name]; !ok {
		databases[name] = database
	}
	//panic(fmt.Sprintf("database %s already registered", name))
}

//Open return opened dababase
func Open(name, connStr string) (*sql.DB, error) {
	if database, ok := databases[name]; ok {
		return database.Open(connStr)
	}
	panic(fmt.Sprintf("not implement database %s", name))
}
