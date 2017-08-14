package policy

import (
	"database/sql"
	"fmt"

	"time"

	"github.com/bocheninc/L0/rest/model/table"
)

func init() {
	Policy := NewPolicy()
	table.Register(Policy.TableName(), Policy)
}

type Policy struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Descr   string    `json:"descr"`
	API     string    `json:"api"`
	Action  string    `json:"action"`
	Params  string    `json:"params"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

//Condition
func (policy *Policy) Condition() (condition string) {
	if policy.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", policy.ID)
	}

	if policy.Name != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" name='%s' ", policy.Name)
	}

	if policy.Descr != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" descr='%s' ", policy.Descr)
	}

	if policy.API != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" api='%s' ", policy.API)
	}

	if policy.Action != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" action='%s' ", policy.Action)
	}

	if policy.Params != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" params='%s' ", policy.Params)
	}
	return
}

//TableName return table name
func (policy *Policy) TableName() string {
	return "policy"
}
func (policy *Policy) validate(tx *sql.Tx) error {
	if policy.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if policy.API == "" {
		return fmt.Errorf("api is empty")
	}
	if policy.Action == "" {
		return fmt.Errorf("action is empty")
	}
	return nil
}

//CreateIfNotExist
func (policy *Policy) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL,
	name VARCHAR(255) NOT NULL UNIQUE,
	descr TEXT,
	api VARCHAR(255) NOT NULL,
	action VARCHAR(255) NOT NULL,
	params TEXT,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, policy.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (policy *Policy) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, name, descr, api, action, params, created, updated from %s", policy.TableName())
	cond := policy.Condition() + condition
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
		policy := NewPolicy()
		if err := rows.Scan(&policy.ID, &policy.Name, &policy.Descr, &policy.API, &policy.Action, &policy.Params, &policy.Created, &policy.Updated); err != nil {
			return res, err
		}
		res = append(res, policy)
	}
	return res, nil
}

//QueryRow
func (policy *Policy) QueryRow(tx *sql.Tx) error {
	row := tx.QueryRow(fmt.Sprintf("select name, descr, api, action, params, created, updated from %s where id=?", policy.TableName()), policy.ID)
	if err := row.Scan(&policy.Name, &policy.Descr, &policy.API, &policy.Action, &policy.Params, &policy.Created, &policy.Updated); err != nil {
		return err
	}
	return nil
}

//Insert
func (policy *Policy) Insert(tx *sql.Tx) error {
	if err := policy.validate(tx); err != nil {
		return err
	}
	policy.Created = time.Now()
	policy.Updated = policy.Created
	policy.ID = policy.getNextID(tx)
	_, err := tx.Exec(fmt.Sprintf("insert into %s(id, name, descr, api, action, params,created, updated) values(?, ?, ?, ?, ?, ?, ?, ?)", policy.TableName()),
		policy.ID, policy.Name, policy.Descr, policy.API, policy.Action, policy.Params, policy.Created, policy.Updated)
	if err != nil {
		return err
	}
	return nil
}

//Delete
func (policy *Policy) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", policy.TableName())
	cond := policy.Condition() + condition
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
func (policy *Policy) Update(tx *sql.Tx) error {
	if err := policy.validate(tx); err != nil {
		return err
	}
	policy.Updated = time.Now()
	res, err := tx.Exec(fmt.Sprintf("update %s set name=?, descr=?, api=?, action=?, params=?, created=?, updated=? where id=?", policy.TableName()),
		policy.Name, policy.Descr, policy.API, policy.Action, policy.Params, policy.Created, policy.Updated, policy.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (policy *Policy) InsertOrUpdate(tx *sql.Tx) error {
	if err := policy.validate(tx); err != nil {
		return err
	}
	policy.Created = time.Now()
	policy.Updated = policy.Created
	policy.ID = policy.getNextID(tx)
	_, err := tx.Exec(fmt.Sprintf("insert into %s(id, name, descr, api, action, params,created, updated) values (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE descr=values(descr), api=values(api), action=values(action), params=values(params), updated=values(updated)", policy.TableName()),
		policy.ID, policy.Name, policy.Descr, policy.API, policy.Action, policy.Params, policy.Created, policy.Updated)
	if err != nil {
		return err
	}
	return nil
}

func (policy *Policy) getNextID(tx *sql.Tx) int64 {
	var id int64 = 1
	row := tx.QueryRow(fmt.Sprintf("select id from %s where id = (select max(id) from %s)", policy.TableName(), policy.TableName()))
	if err := row.Scan(&id); err == nil {
		id = id * 2
	}
	return id
}

//NewPolicy return a Policy object
func NewPolicy() *Policy {
	return &Policy{}
}
