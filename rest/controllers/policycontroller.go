package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/policy"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type PolicyController struct {
}

func (policyCtrl *PolicyController) Get(c *gin.Context) {
	policy := policy.NewPolicy()
	if err := c.BindJSON(&policy); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	policys, err := policy.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, policys)
}

func (policyCtrl *PolicyController) Post(c *gin.Context) {
	policy := policy.NewPolicy()
	if err := c.BindJSON(policy); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := policy.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, policy)
}

func (policyCtrl *PolicyController) Put(c *gin.Context) {
	policy := policy.NewPolicy()
	if err := c.BindJSON(policy); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := policy.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, policy)
}

func (policyCtrl *PolicyController) Delete(c *gin.Context) {
	policy := policy.NewPolicy()
	if err := c.BindJSON(policy); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, _ := model.DB.Begin()
	if err := policy.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, policy)
}

func NewPolicyController() *PolicyController {
	return &PolicyController{}
}
