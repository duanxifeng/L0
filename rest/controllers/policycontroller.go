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

func (policyCtrl *PolicyController) MultipleGet(c *gin.Context) {
	type MultiplePolicy struct {
		PolicyID int64 `json:"policy_id"`
	}
	mPolicy := &MultiplePolicy{}
	if err := c.BindJSON(&mPolicy); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tpolicy := policy.NewPolicy()
	policys, err := tpolicy.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	var tpolicys []*policy.Policy
	for _, tpolicy := range policys {
		if mPolicy.PolicyID&tpolicy.(*policy.Policy).ID > 0 {
			tpolicys = append(tpolicys, tpolicy.(*policy.Policy))
		}
	}
	c.JSON(http.StatusOK, tpolicys)
}

func NewPolicyController() *PolicyController {
	return &PolicyController{}
}
