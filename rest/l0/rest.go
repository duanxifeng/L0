package L0

import (
	"time"

	"github.com/bocheninc/L0/core/types"
	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/transaction"
)

func AddBlock(block types.Block) error {
	txs, err := block.GetTransactions(100)
	if err != nil {
		return err
	}
	transactions := make([]*transaction.Transaction, 0)
	for _, tx := range txs {
		transactions = append(transactions, &transaction.Transaction{
			Hash:      tx.Hash().String(),
			FromChain: tx.FromChain(),
			ToChain:   tx.ToChain(),
			Type:      int64(tx.GetType()),
			Nonce:     int64(tx.Nonce()),
			Sender:    tx.Sender().String(),
			Receiver:  tx.Recipient().String(),
			Amount:    tx.Amount().Uint64(),
			Fee:       tx.Fee().Uint64(),
			Created:   time.Unix(int64(tx.CreateTime()), 0),
		})
	}

	tx, _ := model.DB.Begin()
	for _, transaction := range transactions {
		if err := transaction.Insert(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil

}
