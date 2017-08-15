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
	FromChain string    `json:"from_chain"`
	ToChain   string    `json:"to_chain"`
	Type      int64     `json:"tx_type"`
	Nonce     int64     `json:"tx_nonce"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Amount    uint64    `json:"amount"`
	Fee       uint64    `json:"fee"`
	Signature string    `json:"signature"`
	Created   time.Time `json:"created"`
	Payload   string    `json:"payload"`
	Hash      string    `json:"hash"`
	Height    uint64    `json:"height"`
}

//Condition
func (transcation *Transaction) Condition() (condition string) {
	if transcation.ID != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" id='%d' ", transcation.ID)
	}

	if transcation.FromChain != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" from_chain='%s' ", transcation.FromChain)
	}

	if transcation.ToChain != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" to_chain='%s' ", transcation.ToChain)
	}

	if transcation.Type != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" type='%d' ", transcation.Type)
	}

	if transcation.Sender != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" sender='%s' ", transcation.Sender)
	}

	if transcation.Receiver != "" {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" receiver='%s' ", transcation.Receiver)
	}

	if transcation.Height != 0 {
		if condition != "" {
			condition += " and "
		}
		condition += fmt.Sprintf(" height='%d' ", transcation.Height)
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
	from_chain VARCHAR(255) NOT NULL,
	to_chain VARCHAR(255) NOT NULL,
	tx_type int NOT NULL,
	tx_nonce int NOT NULL,
	sender VARCHAR(255) NOT NULL,
	receiver VARCHAR(255) NOT NULL,
	amount INT NOT NULL,
	fee INT NOT NULL,
	signature VARCHAR(255) NOT NULL,
	created DATETIME NOT NULL,
	payload TEXT NOT NULL,
	hash VARCHAR(255) NOT NULL UNIQUE,
	height INT NOT NULL DEFAULT 0,
	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`
	sql = fmt.Sprintf(sql, transcation.TableName())
	_, err := db.Exec(sql)
	return sql, err
}

//Query
func (transcation *Transaction) Query(db *sql.DB, condition string) ([]table.ITable, error) {
	sql := fmt.Sprintf("select id, from_chain, to_chain, tx_type, tx_nonce, sender, receiver, amount, fee, signature, created, payload, hash, height from %s", transcation.TableName())
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
		if err := rows.Scan(&tx.ID, &tx.FromChain, &tx.ToChain, &tx.Type, &tx.Nonce, &tx.Sender, &tx.Receiver, &tx.Amount, &tx.Fee, &tx.Signature, &tx.Created, &tx.Payload, &tx.Hash, &tx.Height); err != nil {
			return res, err
		}
		res = append(res, tx)
	}
	return res, nil
}

//Insert
func (transcation *Transaction) Insert(tx *sql.Tx) error {
	transcation.Created = time.Now()
	res, err := tx.Exec(fmt.Sprintf("insert into %s(from_chain, to_chain, tx_type, tx_nonce, sender, receiver, amount, fee, signature, created, payload, hash, height) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", transcation.TableName()),
		transcation.FromChain, transcation.ToChain, transcation.Type, transcation.Nonce,
		transcation.Sender, transcation.Receiver, transcation.Amount, transcation.Fee, transcation.Signature, transcation.Created, transcation.Payload, transcation.Hash, transcation.Height)
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
	res, err := tx.Exec(fmt.Sprintf("update %s set from_chain=?, to_chain=?, tx_type=?, tx_nonce=?, sender=?, receiver=?, amount=?, fee=?, signature=?, created=?, payload=? ,hash=?, height=? where id=? ", transcation.TableName()),
		transcation.FromChain, transcation.ToChain, transcation.Type, transcation.Nonce,
		transcation.Sender, transcation.Receiver, transcation.Amount, transcation.Fee, transcation.Signature, transcation.Created, transcation.Payload, transcation.Hash, transcation.Height, transcation.ID)
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
