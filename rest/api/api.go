package api

import (
	"encoding/json"
	"strings"
	"time"

	gin "gopkg.in/gin-gonic/gin.v1"

	"fmt"

	"github.com/bocheninc/L0/rest/controllers"
	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/account"
	"github.com/bocheninc/L0/rest/model/table/history"
	"github.com/bocheninc/L0/rest/model/table/policy"
	"github.com/bocheninc/L0/rest/model/table/status"
	"github.com/bocheninc/L0/rest/model/table/user"
	"github.com/bocheninc/L0/rest/router"
)

//init user table data
var users = make([]*user.User, 0)

//init status table data
var statuses = make([]*status.Status, 0)

//init policy table data
var policys = make([]*policy.Policy, 0)
var handlers = make([]func(c *gin.Context), 0)

//Run start api service
func Run(addr ...string) {
	// create table data
	for _, status := range statuses {
		res, err := status.Query(model.DB, "")
		if err != nil {
			panic(err)
		}
		tx, _ := model.DB.Begin()
		if len(res) == 0 {
			err = status.InsertOrUpdate(tx)
		}
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		tx.Commit()
	}

	for _, opolicy := range policys {
		res, err := opolicy.Query(model.DB, "")
		if err != nil {
			panic(err)
		}
		tx, _ := model.DB.Begin()
		if len(res) == 0 {
			err = opolicy.InsertOrUpdate(tx)
		}
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		tx.Commit()
	}

	statusID := int64(1)
	policyID := int64(0)
	tpolicy := &policy.Policy{}
	res, err := tpolicy.Query(model.DB, "")
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		policyID = r.(*policy.Policy).ID | policyID
	}
	for _, user := range users {
		user.StatusID = statusID
		user.PolicyID = policyID
		res, err := user.Query(model.DB, "")
		if err != nil {
			panic(err)
		}
		tx, _ := model.DB.Begin()
		if len(res) == 0 {
			err = user.InsertOrUpdate(tx)
		}
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		tx.Commit()
	}

	// router register
	for index, opolicy := range policys {
		function := handlers[index]
		var params map[string]string
		json.Unmarshal([]byte(opolicy.Params), &params)
		action := opolicy.Action
		name := opolicy.Name
		router.RegisterPost(opolicy.API, func(c *gin.Context) {
			url := c.Request.URL.String()
			if !strings.Contains(url, "-get") {
				for k := range params {
					params[k] = c.PostForm(k)
				}
				paramsBytes, _ := json.Marshal(params)
				// add history
				history := &history.History{
					API:     url,
					Action:  action,
					Name:    name,
					Params:  string(paramsBytes),
					Created: time.Now(),
					Updated: time.Now(),
				}
				tx, _ := model.DB.Begin()
				if err := history.Insert(tx); err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}
			function(c)
		})
	}

	router.Run(addr...)
}

func init() {
	users = append(users, &user.User{
		Name:     "admin",
		PassWord: "admin",
	})

	statuses = append(statuses, &status.Status{
		Name:  "待审批",
		Descr: "等待审批中...",
	})
	statuses = append(statuses, &status.Status{
		Name:  "失败",
		Descr: "未通过审核",
	})
	statuses = append(statuses, &status.Status{
		Name:  "通过",
		Descr: "已通过审核",
	})
	statuses = append(statuses, &status.Status{
		Name:  "注销",
		Descr: "已销户,不可用",
	})

	//***********************************************************************************
	//全局信息 --- 可操作权限
	//***********************************************************************************
	{
		action := "全局信息"
		opolicy := &policy.Policy{
		// ID         int64         `json:"id"`
		// Name       string        `json:"name"`
		// Descr      string        `json:"descr"`
		// API        string        `json:"api"`
		// Action     string        `json:"action"`
		// Params     string        `json:"params"`
		// Created    time.Time     `json:"created"`
		// Updated    time.Time     `json:"updated"`
		}
		bytes, _ := json.Marshal(opolicy)
		policyCtrl := controllers.NewPolicyController()

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:操作ID"
		descrs["name"] = "name:操作名称"
		descrs["descr"] = "descr:操作描述"
		descrs["api"] = "api:操作URL"
		descrs["action"] = "action:操作类型"
		descrs["params"] = "params:操作URL所需参数的模板"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "descr")
		delete(params, "api")
		delete(params, "params")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的用户账户可用的操作信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户可用的操作查询",
			Descr:  descrStr,
			API:    "/policy-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, policyCtrl.Get)
	}

	//***********************************************************************************
	//全局信息 --- 用户账户及账号的状态
	//***********************************************************************************
	{
		action := "全局信息"
		status := &status.Status{
		// ID    int64        `json:"id"`
		// Name    string        `json:"name"`
		// Descr    string        `json:"descr"`
		// Created    time.Time    `json:"created"`
		// Updated    time.Time    `json:"updated"`
		}
		bytes, _ := json.Marshal(status)
		statusCtrl := controllers.NewStatusController()

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:状态ID"
		descrs["name"] = "name:状态名称"
		descrs["descr"] = "descr:状态描述"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "descr")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的用户账号及账号的可用状态信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}

		policys = append(policys, &policy.Policy{
			Name:   "用户账号及账号的可用状态查询",
			Descr:  descrStr,
			API:    "/status-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, statusCtrl.Get)
	}

	//***********************************************************************************
	//全局信息 --- 操作记录
	//***********************************************************************************
	{
		action := "全局信息"
		history := &history.History{
		// ID    int64        `json:"id"`
		// Name    string        `json:"name"`
		// API    string        `json:"api"`
		// Action    string        `json:"action"`
		// Params    string        `json:"params"`
		// Created    time.Time    `json:"created"`
		// Updated    time.Time    `json:"updated"`
		}
		historyCtrl := controllers.NewHistoryController()
		bytes, _ := json.Marshal(history)

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:操作记录ID"
		descrs["name"] = "name:操作记录名称"
		descrs["api"] = "api:操作api"
		descrs["action"] = "name:操作类型"
		descrs["params"] = "params:操作参数"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "api")
		delete(params, "params")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的操作记录信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号的操作记录查询",
			Descr:  descrStr,
			API:    "/history-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, historyCtrl.Get)
	}

	//***********************************************************************************
	//数据处理基本功能 --- 用户账户
	//***********************************************************************************
	{
		action := "数据处理基本功能"
		user := &user.User{
		// ID          int64        `json:"id"`
		// Name        string       `json:"name"`
		// PassWord    string       `json:"password"`
		// Metadata    string       `json:"metadata"`
		// PolicyID    int64        `json:"policy_id"`
		// StatusID    int64        `json:"status_id"`
		// Created     time.Time    `json:"created"`
		// Updated     time.Time    `json:"updated"`
		}
		bytes, _ := json.Marshal(user)
		userCtrl := controllers.NewUserController()

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:用户账户ID"
		descrs["name"] = "name:用户账户名称"
		descrs["password"] = "password:用户账户密码"
		descrs["metadata"] = "metadata:用户账户附加信息"
		descrs["policy_id"] = "policy_id:用户账户权限ID"
		descrs["status_id"] = "status_id:用户账户状态ID"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "password")
		delete(params, "metadata")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的用户账户信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户查询",
			Descr:  descrStr,
			API:    "/user-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, userCtrl.Get)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "updated")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `创建用户账户
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户创建",
			Descr:  descrStr,
			API:    "/user-post",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, userCtrl.Post)

		json.Unmarshal(bytes, &params)
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `修改用户账户信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户修改",
			Descr:  descrStr,
			API:    "/user-put",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, userCtrl.Put)

		json.Unmarshal(bytes, &params)
		delete(params, "name")
		delete(params, "password")
		delete(params, "metadata")
		delete(params, "policy_id")
		delete(params, "status_id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `注销用户账户
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户注销",
			Descr:  descrStr,
			API:    "/user-close",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, userCtrl.Put)
	}

	//***********************************************************************************
	//数据处理基本功能 --- 账号
	//***********************************************************************************
	{
		action := "数据处理基本功能"
		account := &account.Account{
		// ID        int64        `json:"id"`
		// Address   string        `json:"addr"`
		// UserID    int64        `json:"user_id"`
		// StatusID  int64        `json:"status_id"`
		}
		bytes, _ := json.Marshal(account)
		accountCtrl := controllers.NewAccountController()

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:用户账号ID"
		descrs["addr"] = "addr:用户账号地址"
		descrs["user_id"] = "user_id:用户账号所属用户账户ID"
		descrs["status_id"] = "status_id:用户账号的状态ID"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的用户账号信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号查询",
			Descr:  descrStr,
			API:    "/account-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Get)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `创建用户账号
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号创建",
			Descr:  descrStr,
			API:    "/account-post",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Post)

		json.Unmarshal(bytes, &params)
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `审批用户账号
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号审批",
			Descr:  descrStr,
			API:    "/account-put",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Put)

		json.Unmarshal(bytes, &params)
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `导出用户账号
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号导出",
			Descr:  descrStr,
			API:    "/account-export",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Post)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `导入用户账号
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号导入",
			Descr:  descrStr,
			API:    "/account-import",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Post)
	}
}
