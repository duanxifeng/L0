package transaction

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bocheninc/L0/rest/model/table"
)

func init() {
	tx := NewTransaction()
	table.Register(tx.TableName(), tx)
}

type Transaction struct {
	ID        int64     `json:"id"`
	Hash      string    `json:"hash"`
	FromChain string    `json:"from_chain"`
	ToChain   string    `json:"to_chain"`
	Type      int64     `json:"tx_type"`
	Nonce     int64     `json:"tx_nonce"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    uint64    `json:"amount"`
	Fee       uint64    `json:"fee"`
	Created   time.Time `json:"created"`
}

//Condition
func (transcation *Transaction) Condition() (condition string) {
	if transcation.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", transcation.ID)
	}
	return
}

//TableName return table name
func (transcation *Transaction) TableName() string {
	return "transaction"
}

//CreateIfNotExist
func (transcation *Transaction) CreateIfNotExist(db *sql.DB) (string, error) {
	sql := `
	CREATE TABLE IF NOT EXISTS %s (
	id INT NOT NULL AUTO_INCREMENT,
	hash CHAR(255) NOT NULL UNIQUE,
	from_chain CHAR(255) NOT NULL,
	to_chain CHAR(255) NOT NULL,
	tx_type int NOT NULL,
	tx_nonce int NOT NULL,
	sender CHAR(255) NOT NULL,
	receiver CHAR(255) NOT NULL,
	amount INT NOT NULL,
	fee INT NOT NULL,
	created DATETIME NOT NULL,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, transcation.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (transcation *Transaction) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, hash, from_chain, to_chain, tx_type, tx_nonce, sender, receiver, amount, fee, created from %s", transcation.TableName())
	cond := transcation.Condition() + condition
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
		tx := NewTransaction()
		if err := rows.Scan(&tx.ID, &tx.Hash, &tx.FromChain, &tx.ToChain, &tx.Type, &tx.Nonce, &tx.Sender, &tx.Receiver, &tx.Amount, tx.Fee, &tx.Created); err != nil {
			return res, err
		}
		res = append(res, tx)
	}
	return res, nil
}

//Insert
func (transcation *Transaction) Insert(tx *sql.Tx) error {
	res, err := tx.Exec("insert info ?(hash, from_chain, to_chain, tx_type, tx_nonce, sender, receiver, amount, fee, created) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		transcation.TableName(), transcation.Hash, transcation.FromChain, transcation.ToChain, transcation.Type, transcation.Nonce,
		transcation.Sender, transcation.Receiver, transcation.Amount, transcation.Fee, transcation.Created)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	transcation.ID = id
	return nil
}

//Delete
func (transcation *Transaction) Delete(tx *sql.Tx, condition string) error {
	sql := fmt.Sprintf("delete from %s", transcation.TableName())
	cond := transcation.Condition() + condition
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
func (transcation *Transaction) Update(tx *sql.Tx) error {
	res, err := tx.Exec(fmt.Sprintf("update %s set hash=?, from_chain=?, to_chain=?, tx_type=?, tx_nonce=?, sender=?, receiver=?, amount=? fee=?, created=? where id=? ", transcation.TableName()),
		transcation.Hash, transcation.FromChain, transcation.ToChain, transcation.Type, transcation.Nonce,
		transcation.Sender, transcation.Receiver, transcation.Amount, transcation.Fee, transcation.Created, transcation.ID)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

//NewTransaction return a tx object
func NewTransaction() *Transaction {
	return &Transaction{}
}
