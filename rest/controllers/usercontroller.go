package controllers

import (
	"io"
	"net/http"

	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/user"
	gin "gopkg.in/gin-gonic/gin.v1"
)

type UserController struct {
}

func (userCtrl *UserController) Get(c *gin.Context) {
	user := user.NewUser()
	if err := c.BindJSON(&user); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	users, err := user.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (userCtrl *UserController) Post(c *gin.Context) {
	user := user.NewUser()
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	tx, _ := model.DB.Begin()

	if err := user.Insert(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, user)
}

func (userCtrl *UserController) Put(c *gin.Context) {
	user := user.NewUser()
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	tx, _ := model.DB.Begin()

	if err := user.Update(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, user)
}

func (userCtrl *UserController) Delete(c *gin.Context) {
	user := user.NewUser()
	if err := c.BindJSON(user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	tx, _ := model.DB.Begin()

	if err := user.Delete(tx, ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, user)
}

func NewUserController() *UserController {
	return &UserController{}
}
