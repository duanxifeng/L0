package status

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"time"

	"github.com/bocheninc/L0/rest/model/table"
)

func init() {
	status := NewStatus()
	table.Register(status.TableName(), status)
}

type Status struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Descr   string    `json:"descr"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (status *Status) GetDescr() string {
	return fmt.Sprintf(` list %s `, status.TableName())
}
func (status *Status) GetParams() string {
	bytes, _ := json.Marshal(status)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (status *Status) PostDescr() string {
	return fmt.Sprintf(` add %s `, status.TableName())
}
func (status *Status) PostParams() string {
	bytes, _ := json.Marshal(status)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "id")
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (status *Status) PutDescr() string {
	return fmt.Sprintf(` modify %s `, status.TableName())
}
func (status *Status) PutParams() string {
	bytes, _ := json.Marshal(status)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (status *Status) DeleteDescr() string {
	return fmt.Sprintf(` delete %s `, status.TableName())
}
func (status *Status) DeleteParams() string {
	bytes, _ := json.Marshal(status)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}

//Condition
func (status *Status) Condition() (condition string) {
	if status.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", status.ID)
	}

	if status.Name != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" name='%s' ", status.Name)
	}

	if status.Descr != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" descr='%s' ", status.Descr)
	}
	return
}

//TableName return table name
func (status *Status) TableName() string {
	return "status"
}
func (status *Status) validate(tx *sql.Tx) error {
	if status.Name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}

//CreateIfNotExist
func (status *Status) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(20) NOT NULL UNIQUE,
	descr VARCHAR(255),
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, status.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (status *Status) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, name, descr, created, updated from %s", status.TableName())
	cond := status.Condition() + condition
	if cond != "" {
		sql = fmt.Sprintf("%s where %s", sql, cond)
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]table.ITable, 0)
	for rows.Next() {
		status := NewStatus()
		if err := rows.Scan(&status.ID, &status.Name, &status.Descr, &status.Created, &status.Updated); err != nil {
			return res, err
		}
		res = append(res, status)
	}
	return res, nil
}

//QueryRow
func (status *Status) QueryRow(tx *sql.Tx) error {
	row := tx.QueryRow(fmt.Sprintf("select name, descr, created, updated from %s where id=?", status.TableName()), status.ID)
	if err := row.Scan(&status.Name, &status.Descr, &status.Created, &status.Updated); err != nil {
		return err
	}
	return nil
}

//Insert
func (status *Status) Insert(tx *sql.Tx) error {
	if err := status.validate(tx); err != nil {
		return err
	}
	status.Created = time.Now()
	status.Updated = status.Created
	res, err := tx.Exec(fmt.Sprintf("insert into %s(name, descr, created, updated) values(?, ?, ?, ?)", status.TableName()),
		status.Name, status.Descr, status.Created, status.Updated)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	status.ID = id
	return nil
}

//Delete
func (status *Status) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", status.TableName())
	cond := status.Condition() + condition
	if cond != "" {
		sql = fmt.Sprintf("%s where %s", sql, cond)
	}

	res, err := tx.Exec(sql)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

//Update
func (status *Status) Update(tx *sql.Tx) error {
	if err := status.validate(tx); err != nil {
		return err
	}
	status.Updated = time.Now()
	res, err := tx.Exec(fmt.Sprintf("update %s set name=?, descr=?, created=?, updated=? where id=?", status.TableName()),
		status.Name, status.Descr, status.Created, status.Updated, status.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (status *Status) InsertOrUpdate(tx *sql.Tx) error {
	if err := status.validate(tx); err != nil {
		return err
	}
	status.Created = time.Now()
	status.Updated = status.Created
	_, err := tx.Exec(fmt.Sprintf("insert into %s(name, descr, created, updated) values(?, ?, ?, ?) ON DUPLICATE KEY UPDATE descr=values(descr), updated=values(updated)", status.TableName()),
		status.Name, status.Descr, status.Created, status.Updated)
	if err != nil {
		return err
	}
	return nil
}

//NewStatus return a Status object
func NewStatus() *Status {
	return &Status{}
}
