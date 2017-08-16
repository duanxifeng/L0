package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bocheninc/L0/rest/model/table"
	"github.com/bocheninc/L0/rest/model/table/policy"
	"github.com/bocheninc/L0/rest/model/table/status"
)

func init() {
	user := NewUser()
	table.Register(user.TableName(), user)
}

type User struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	PassWord string    `json:"password"`
	Metadata string    `json:"metadata"`
	PolicyID int64     `json:"policy_id"`
	StatusID int64     `json:"status_id"`
	Range    string    `json:"query_range"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	policys  []*policy.Policy
	status   *status.Status
}

//Condition
func (user *User) Condition() (condition string) {
	if user.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", user.ID)
	}

	if user.Name != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" name='%s' ", user.Name)
	}

	if user.PolicyID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" policy_id='%d' ", user.PolicyID)
	}

	if user.StatusID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" status_id='%d' ", user.StatusID)
	}
	return
}

//TableName return table name
func (user *User) TableName() string {
	return "user"
}
func (user *User) validate(tx *sql.Tx) error {
	if user.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if user.PassWord == "" {
		return fmt.Errorf("password is empty")
	}
	if _, err := user.Policys(tx); err != nil {
		return fmt.Errorf("policys_id %d is not exist", user.PolicyID)
	}
	if _, err := user.Status(tx); err != nil {
		return fmt.Errorf("status_id %d is not exist", user.StatusID)
	}
	return nil
}

//CreateIfNotExist
func (user *User) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	metadata TEXT NOT NULL,
	policy_id INT NOT NULL,
	status_id INT NOT NULL,
	query_range VARCHAR(255) NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, user.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (user *User) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, name, password, metadata, policy_id, status_id, query_range, created, updated from %s", user.TableName())
	cond := user.Condition() + condition
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
		user := NewUser()
		if err := rows.Scan(&user.ID, &user.Name, &user.PassWord, &user.Metadata, &user.PolicyID, &user.StatusID, &user.Range, &user.Created, &user.Updated); err != nil {
			return res, err
		}
		res = append(res, user)
	}
	return res, nil
}

//QueryRow
func (user *User) QueryRow(tx *sql.Tx) error {
	row := tx.QueryRow(fmt.Sprintf("select name, password, metadata, policy_id, status_id, query_range, created, updated from %s where id=?", user.TableName()), user.ID)
	if err := row.Scan(&user.Name, &user.PassWord, &user.Metadata, &user.PolicyID, &user.StatusID, &user.Range, &user.Created, &user.Updated); err != nil {
		return err
	}
	return nil
}

//Insert
func (user *User) Insert(tx *sql.Tx) error {
	if err := user.validate(tx); err != nil {
		return err
	}
	user.Created = time.Now()
	user.Updated = user.Created
	res, err := tx.Exec(fmt.Sprintf("insert into %s(name, password, metadata, policy_id, status_id, query_range, created, updated) values(?, ?, ?, ?, ?, ?, ?, ?)", user.TableName()),
		user.Name, user.PassWord, user.Metadata, user.PolicyID, user.StatusID, user.Range, user.Created, user.Updated)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

//Delete
func (user *User) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", user.TableName())
	cond := user.Condition() + condition
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
func (user *User) Update(tx *sql.Tx) error {
	if err := user.validate(tx); err != nil {
		return err
	}
	user.Updated = time.Now()
	res, err := tx.Exec(fmt.Sprintf("update %s set name=?, password=?, metadata=?, policy_id=?, status_id=?, query_range, created=?, updated=? where id=? ", user.TableName()),
		user.Name, user.PassWord, user.Metadata, user.PolicyID, user.StatusID, user.Range, user.Created, user.Updated, user.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (user *User) InsertOrUpdate(tx *sql.Tx) error {
	if err := user.validate(tx); err != nil {
		return err
	}
	user.Created = time.Now()
	user.Updated = user.Created
	_, err := tx.Exec(fmt.Sprintf("insert into %s(name, password, metadata, policy_id, status_id, query_range, created, updated) values(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE password=values(password), metadata=values(metadata), policy_id=values(policy_id), status_id= values(status_id), query_range=values(query_range), created=values(created), updated=values(updated)", user.TableName()),
		user.Name, user.PassWord, user.Metadata, user.PolicyID, user.StatusID, user.Range, user.Created, user.Updated)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) Policys(tx *sql.Tx) ([]*policy.Policy, error) {
	if len(user.policys) == 0 {
		policys := make([]*policy.Policy, 0)

		oper := int64(1)
		for i := uint16(1); ; i++ {
			if tid := user.PolicyID & oper; tid > 0 {
				policy := policy.NewPolicy()
				policy.ID = tid
				if err := policy.QueryRow(tx); err != nil {
					return nil, err
				}
				policys = append(policys)
			}
			oper = oper << 1
			if user.PolicyID>>i == 0 {
				break
			}
		}
		user.policys = policys
	}
	return user.policys, nil
}

func (user *User) Status(tx *sql.Tx) (*status.Status, error) {
	if user.status == nil {
		user.status = status.NewStatus()
		user.status.ID = user.StatusID
		if err := user.status.QueryRow(tx); err != nil {
			return nil, err
		}
	}
	return user.status, nil
}

//NewUser return a user object
func NewUser() *User {
	return &User{}
}
