package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"fmt"

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
	var values map[string]interface{}
	if err := c.BindJSON(&values); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := user.NewUser()
	user.ID = int64(values["id"].(float64))
	tx, _ := model.DB.Begin()
	if err := user.QueryRow(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	bytes, _ := json.Marshal(user)
	var tvalues map[string]interface{}
	json.Unmarshal(bytes, &tvalues)
	for k, v := range values {
		tvalues[k] = v
	}
	bytes, _ = json.Marshal(tvalues)
	json.Unmarshal(bytes, user)

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

func (userCtrl *UserController) Login(c *gin.Context) {
	tuser := user.NewUser()
	if err := c.BindJSON(&tuser); err != nil && err != io.EOF {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	users, err := tuser.Query(model.DB, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Errorf("用户名不存在").Error(),
		})
		return
	}

	if users[0].(*user.User).PassWord != tuser.PassWord {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Errorf("用户名与密码验证失败").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users[0])
}

func NewUserController() *UserController {
	return &UserController{}
}
