package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/history"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type HistoryController struct {
}

func (historyCtrl *HistoryController) Get(c *gin.Context) {
	history := history.NewHistory()
	if err := c.BindJSON(&history); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	historys, err := history.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, historys)
}

func (historyCtrl *HistoryController) Post(c *gin.Context) {
	history := history.NewHistory()
	if err := c.BindJSON(history); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := history.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, history)
}

func (historyCtrl *HistoryController) Put(c *gin.Context) {
	history := history.NewHistory()
	if err := c.BindJSON(history); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := history.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, history)
}

func (historyCtrl *HistoryController) Delete(c *gin.Context) {
	history := history.NewHistory()
	if err := c.BindJSON(history); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := history.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, history)
}

func NewHistoryController() *HistoryController {
	return &HistoryController{}
}
