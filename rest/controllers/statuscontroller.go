package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/status"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type StatusController struct {
}

func (statusCtrl *StatusController) Get(c *gin.Context) {
	status := status.NewStatus()
	if err := c.BindJSON(&status); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	statuses, err := status.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (statusCtrl *StatusController) Post(c *gin.Context) {
	status := status.NewStatus()
	if err := c.BindJSON(status); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := status.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, status)
}

func (statusCtrl *StatusController) Put(c *gin.Context) {
	status := status.NewStatus()
	if err := c.BindJSON(status); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := status.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, status)
}

func (statusCtrl *StatusController) Delete(c *gin.Context) {
	status := status.NewStatus()
	if err := c.BindJSON(status); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := status.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, status)
}

func NewStatusController() *StatusController {
	return &StatusController{}
}
