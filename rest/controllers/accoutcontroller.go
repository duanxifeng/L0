package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/account"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type AccountController struct {
}

func (accountCtrl *AccountController) Get(c *gin.Context) {
	account := account.NewAccount()
	if err := c.BindJSON(&account); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	accounts, err := account.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (accountCtrl *AccountController) Post(c *gin.Context) {
	account := account.NewAccount()
	if err := c.BindJSON(account); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := account.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, account)
}

func (accountCtrl *AccountController) Put(c *gin.Context) {
	account := account.NewAccount()
	if err := c.BindJSON(account); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := account.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, account)
}

func (accountCtrl *AccountController) Delete(c *gin.Context) {
	account := account.NewAccount()
	if err := c.BindJSON(account); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := account.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, account)
}

func NewAccountController() *AccountController {
	return &AccountController{}
}
