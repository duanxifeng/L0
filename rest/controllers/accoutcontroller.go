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

type ExportObj struct {
	User    string `json:"user"`
	Account string `json:"account"`
}

func (accountCtrl *AccountController) Export(c *gin.Context) {
	taccount := account.NewAccount()
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

	objs := make([]*ExportObj, 0)
	tx, _ := model.DB.Begin()
	for _, taccount := range accounts {
		user, err := taccount.(*account.Account).User(tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		objs = append(objs, &ExportObj{
			User:    user.Name,
			Account: taccount.(*account.Account).Address,
		})
	}
	tx.Commit()

	c.JSON(http.StatusOK, objs)
}

func (accountCtrl *AccountController) Import(c *gin.Context) {
	var objs []*ExportObj
	if err := c.BindJSON(&objs); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	accounts := make([]*account.Account, 0)
	for _, obj := range objs {
		taccount := &account.Account{
			Address:  obj.Account,
			UserID:   3,
			StatusID: 5,
		}
		accounts = append(accounts, taccount)
	}

	c.JSON(http.StatusOK, accounts)
}

func NewAccountController() *AccountController {
	return &AccountController{}
}
