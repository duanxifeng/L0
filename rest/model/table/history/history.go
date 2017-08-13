package history

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bocheninc/L0/rest/model/table"
)

func init() {
	history := NewHistory()
	table.Register(history.TableName(), history)
}

type History struct {
	ID      int64     `json:"id"`
	API     string    `json:"api"`
	Action  string    `json:"action"`
	Params  string    `json:"params"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (history *History) GetDescr() string {
	return fmt.Sprintf(` list %s `, history.TableName())
}
func (history *History) GetParams() string {
	bytes, _ := json.Marshal(history)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (history *History) PostDescr() string {
	return fmt.Sprintf(` add %s `, history.TableName())
}
func (history *History) PostParams() string {
	bytes, _ := json.Marshal(history)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "id")
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (history *History) PutDescr() string {
	return fmt.Sprintf(` modify %s `, history.TableName())
}
func (history *History) PutParams() string {
	bytes, _ := json.Marshal(history)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (history *History) DeleteDescr() string {
	return fmt.Sprintf(` delete %s `, history.TableName())
}
func (history *History) DeleteParams() string {
	bytes, _ := json.Marshal(history)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}

//Condition
func (history *History) Condition() (condition string) {
	if history.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", history.ID)
	}

	if history.API != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" api='%s' ", history.API)
	}

	if history.Action != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" action='%s' ", history.Action)
	}

	if history.Params != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" params='%s' ", history.Action)
	}
	return
}

//TableName return table name
func (history *History) TableName() string {
	return "history"
}
func (history *History) validate(tx *sql.Tx) error {
	return nil
}

//CreateIfNotExist
func (history *History) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL AUTO_INCREMENT,
	api VARCHAR(255) NOT NULL,
	action VARCHAR(255) NOT NULL,
	params TEXT,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, history.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (history *History) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, api, action, params, created, updated from %s", history.TableName())
	cond := history.Condition() + condition
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
		history := NewHistory()
		if err := rows.Scan(&history.ID, &history.API, &history.Action, &history.Params, &history.Created, &history.Updated); err != nil {
			return res, err
		}
		res = append(res, history)
	}
	return res, nil
}

//Insert
func (history *History) Insert(tx *sql.Tx) error {
	if err := history.validate(tx); err != nil {
		return err
	}
	history.Created = time.Now()
	history.Updated = history.Created
	res, err := tx.Exec(fmt.Sprintf("insert into %s(api, action, params, created, updated) values(?, ?, ?, ?, ?)", history.TableName()),
		history.API, history.Action, history.Params, history.Created, history.Updated)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	history.ID = id
	return nil
}

//Delete
func (history *History) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", history.TableName())
	cond := history.Condition() + condition
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
func (history *History) Update(tx *sql.Tx) error {
	if err := history.validate(tx); err != nil {
		return err
	}
	history.Updated = time.Now()
	res, err := tx.Exec(fmt.Sprintf("update %s set api=?, action=?, params=?,  created=?, updated=? where id=? ", history.TableName()),
		history.API, history.Action, history.Params, history.Created, history.Updated, history.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

//NewHistory return a history object
func NewHistory() *History {
	return &History{}
}
