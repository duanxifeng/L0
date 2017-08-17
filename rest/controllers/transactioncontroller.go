package controllers

import (
	"io"
	"net/http"

	"fmt"

	"strings"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/account"
	"github.com/bocheninc/L0/rest/model/table/transaction"
	"github.com/bocheninc/L0/rest/model/table/user"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type TransactionController struct {
}

func (transactionCtrl *TransactionController) Get(c *gin.Context) {
	transaction := transaction.NewTransaction()
	if err := c.BindJSON(&transaction); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	transactions, err := transaction.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (transactionCtrl *TransactionController) Post(c *gin.Context) {
	transaction := transaction.NewTransaction()
	if err := c.BindJSON(transaction); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := transaction.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, transaction)
}

func (transactionCtrl *TransactionController) Put(c *gin.Context) {
	transaction := transaction.NewTransaction()
	if err := c.BindJSON(transaction); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := transaction.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, transaction)
}

func (transactionCtrl *TransactionController) Delete(c *gin.Context) {
	transaction := transaction.NewTransaction()
	if err := c.BindJSON(transaction); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := transaction.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, transaction)
}

func (transactionCtrl *TransactionController) History(c *gin.Context) {
	taccount := &account.Account{}
	if err := c.BindJSON(&taccount); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	accounts, err := taccount.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Errorf("addr or user not found").Error(),
		})
		return
	}

	senders := " sender in("
	receivers := " receiver in("
	for index, taccount := range accounts {
		addr := taccount.(*account.Account).Address
		if index != 0 {
			senders += ","
			receivers += ","
		}
		senders += fmt.Sprintf("'%s'", addr)
		receivers += fmt.Sprintf("'%s'", addr)
	}
	senders += ")"
	receivers += ")"

	transaction := transaction.NewTransaction()
	transactions, err := transaction.Query(model.DB, senders+" or "+receivers)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (transactionCtrl *TransactionController) QueryRange(c *gin.Context) {
	tuser := &user.User{}
	if err := c.BindJSON(&tuser); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := tuser.QueryRow(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	tx.Commit()

	userIDs := strings.Split(tuser.Range, ",")
	ids := " user_id in("
	for index, userID := range userIDs {
		if index != 0 {
			ids += ","
		}
		ids += fmt.Sprintf("'%s'", userID)
	}
	ids += ")"

	taccount := &account.Account{}
	accounts, err := taccount.Query(model.DB, ids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Errorf("addr or user not found").Error(),
		})
		return
	}

	senders := " sender in("
	receivers := " receiver in("
	for index, taccount := range accounts {
		addr := taccount.(*account.Account).Address
		if index != 0 {
			senders += ","
			receivers += ","
		}
		senders += fmt.Sprintf("'%s'", addr)
		receivers += fmt.Sprintf("'%s'", addr)
	}
	senders += ")"
	receivers += ")"

	transaction := transaction.NewTransaction()
	transactions, err := transaction.Query(model.DB, senders+" or "+receivers)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}
