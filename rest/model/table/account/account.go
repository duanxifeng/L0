package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bocheninc/L0/rest/model/table"
	"github.com/bocheninc/L0/rest/model/table/user"
)

func init() {
	account := NewAccount()
	table.Register(account.TableName(), account)
}

type Account struct {
	ID      int64     `json:"id"`
	Address string    `json:"addr"`
	UserID  int64     `json:"user_id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	user    *user.User
}

func (account *Account) GetDescr() string {
	return fmt.Sprintf(` list %s `, account.TableName())
}
func (account *Account) GetParams() string {
	bytes, _ := json.Marshal(account)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (account *Account) PostDescr() string {
	return fmt.Sprintf(` add %s `, account.TableName())
}
func (account *Account) PostParams() string {
	bytes, _ := json.Marshal(account)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "id")
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (account *Account) PutDescr() string {
	return fmt.Sprintf(` modify %s `, account.TableName())
}
func (account *Account) PutParams() string {
	bytes, _ := json.Marshal(account)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}
func (account *Account) DeleteDescr() string {
	return fmt.Sprintf(` delete %s `, account.TableName())
}
func (account *Account) DeleteParams() string {
	bytes, _ := json.Marshal(account)
	var params map[string]interface{}
	json.Unmarshal(bytes, &params)
	delete(params, "created")
	delete(params, "updated")
	bytes, _ = json.Marshal(params)
	return string(bytes)
}

//Condition
func (account *Account) Condition() (condition string) {
	if account.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", account.ID)
	}

	if account.Address != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" addr='%s' ", account.Address)
	}

	if account.UserID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" user_id='%d' ", account.UserID)
	}
	return
}

//TableName return table name
func (account *Account) TableName() string {
	return "account"
}

//CreateIfNotExist
func (account *Account) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL AUTO_INCREMENT,
	addr VARCHAR(400) NOT NULL,
	user_id INT NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, account.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (account *Account) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, addr, user_id, created, updated from %s", account.TableName())
	cond := account.Condition() + condition
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
		account := NewAccount()
		if err := rows.Scan(&account.ID, &account.Address, &account.UserID, &account.Created, &account.Updated); err != nil {
			return res, err
		}
		res = append(res, account)
	}
	return res, nil
}

//Insert
func (account *Account) Insert(tx *sql.Tx) error {
	account.Created = time.Now()
	account.Updated = account.Created
	res, err := tx.Exec(fmt.Sprintf("insert into %s(addr, user_id, created, updated) values(?, ?, ?, ?)", account.TableName()),
		account.Address, account.UserID, account.Created, account.Updated)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	account.ID = id
	return nil
}

//Delete
func (account *Account) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", account.TableName())
	cond := account.Condition() + condition
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
func (account *Account) Update(tx *sql.Tx) error {
	account.Updated = time.Now()
	res, err := tx.Exec(fmt.Sprintf("update %s set addr=?, user_id=?,  created=?, updated=? where id=? ", account.TableName()),
		account.Address, account.UserID, account.Created, account.Updated, account.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (account *Account) User(tx *sql.Tx) (*user.User, error) {
	if account.user == nil {
		account.user = user.NewUser()
		account.user.ID = account.UserID
		if err := account.user.QueryRow(tx); err != nil {
			return nil, err
		}
	}
	return account.user, nil
}

//NewAccout return a account object
func NewAccount() *Account {
	return &Account{}
}
