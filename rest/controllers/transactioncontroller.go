package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/transaction"
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

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}
