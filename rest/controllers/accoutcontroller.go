package controllers

import (
	"encoding/json"
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
	var values map[string]interface{}
	if err := c.BindJSON(&values); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	account := account.NewAccount()
	account.ID = int64(values["id"].(float64))
	tx, _ := model.DB.Begin()
	if err := account.QueryRow(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	bytes, _ := json.Marshal(account)
	var tvalues map[string]interface{}
	json.Unmarshal(bytes, &tvalues)
	for k, v := range values {
		tvalues[k] = v
	}
	bytes, _ = json.Marshal(tvalues)
	json.Unmarshal(bytes, account)

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

type ExportAccount struct {
	Address string `json:"addr"`
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

	objs := make([]*ExportAccount, 0)
	for _, taccount := range accounts {
		objs = append(objs, &ExportAccount{
			Address: taccount.(*account.Account).Address,
		})
	}

	c.JSON(http.StatusOK, objs)
}

func (accountCtrl *AccountController) Import(c *gin.Context) {
	var accounts []*account.Account
	if err := c.BindJSON(&accounts); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	res := make(map[string]string, 0)
	for _, account := range accounts {
		account.StatusID = 1
		tx, _ := model.DB.Begin()
		if err := account.Insert(tx); err != nil {
			tx.Rollback()
			res[account.Address] = err.Error()
			continue
		}
		bytes, _ := json.Marshal(account)
		res[account.Address] = string(bytes)
		tx.Commit()
	}
	c.JSON(http.StatusOK, res)
}

func (accountCtrl *AccountController) UserInfo(c *gin.Context) {
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

func NewAccountController() *AccountController {
	return &AccountController{}
}
