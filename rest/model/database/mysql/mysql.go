package mysql

import (
	"database/sql"

	"github.com/bocheninc/L0/rest/model/database"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	mysql := NewMySQL()
	database.Register(mysql.Name(), mysql)
}

//MySQL  mysql
type MySQL struct {
}

//Name return database driver name
func (mysql *MySQL) Name() string {
	return "mysql"
}

//Open return opened database
func (mysql *MySQL) Open(connStr string) (*sql.DB, error) {
	return sql.Open(mysql.Name(), connStr)
}

//NewMySQL return a MySQL object
func NewMySQL() *MySQL {
	return &MySQL{}
}
