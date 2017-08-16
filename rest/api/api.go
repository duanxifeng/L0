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
	"github.com/bocheninc/L0/rest/model/table/transaction"
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
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Content-type", "application/json")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			function(c)
		})
	}

	router.Run(addr...)
}

func init() {
	users = append(users, &user.User{
		Name:     "admin",
		PassWord: "admin",
		//Metadata: ""
		StatusID: 3, //已通过审核
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
		descrs["params"] = "params:操作URL所需参数"
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

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
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

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
		descrs["action"] = "action:操作类型"
		descrs["params"] = "params:操作参数"
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

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
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

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
		params["policy_id"] = 0
		params["status_id"] = 1
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
		delete(params, "policy_id")
		delete(params, "status_id")
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
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `审批用户账户
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户审批",
			Descr:  descrStr,
			API:    "/user-status",
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
		// PassWord  string       `json:"password"`
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
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

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
		delete(params, "user_id")
		delete(params, "created")
		delete(params, "updated")
		params["status_id"] = 1
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
		delete(params, "addr")
		delete(params, "status_id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `用户账号映射用户账户
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账号映射用户账户",
			Descr:  descrStr,
			API:    "/account-put",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, accountCtrl.Put)

		json.Unmarshal(bytes, &params)
		delete(params, "addr")
		delete(params, "user_id")
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
			API:    "/account-status",
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
		handlers = append(handlers, accountCtrl.Export)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "user_id")
		delete(params, "status_id")
		delete(params, "created")
		delete(params, "updated")
		var p []map[string]interface{}
		paramsBytes, _ = json.Marshal(append(p, params))
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
		handlers = append(handlers, accountCtrl.Import)
	}

	//***********************************************************************************
	//数据处理基本功能 --- 转账
	//***********************************************************************************
	{
		action := "数据处理基本功能"
		transaction := &transaction.Transaction{
		// ID        int64        `json:"id"`
		// FromChain    string        `json:"from_chain"`
		// ToChain        string        `json:"to_chain"`
		// Type        int64        `json:"tx_type"`
		// Nonce        int64        `json:"tx_nonce"`
		// Sender        string        `json:"sender"`
		// Receiver    string        `json:"receiver"`
		// Amount        uint64        `json:"amount"`
		// Fee        uint64        `json:"fee"`
		// Signature    string        `json:"signature"`
		// Created        time.Time    `json:"created"`
		// Payload        string        `json:"payload"`
		// Hash        string        `json:"hash"`
		// Height        uint64        `json:"height"`
		}
		bytes, _ := json.Marshal(transaction)
		transactionCtrl := controllers.NewTransactionController()

		descrs := make(map[string]string, 0)
		descrs["id"] = "id:交易ID"
		descrs["from_chain"] = "from_chain:交易来源Chain"
		descrs["to_chain"] = "to_chain:交易目的Chain"
		descrs["tx_type"] = "tx_type:交易类型"
		descrs["tx_nonce"] = "tx_nonce:交易Nonce"
		descrs["sender"] = "sender:交易发送方,如果多个,以逗号区分"
		descrs["receiver"] = "receiver:交易接收方,如果多个,以逗号区分"
		descrs["amount"] = "amount:交易金额"
		descrs["fee"] = "fee:更新手续费"
		descrs["signature"] = "signature:交易签名"
		descrs["created"] = "created:交易时间"
		descrs["payload"] = "payload:交易附加信息"
		descrs["hash"] = "hash:交易哈希值"
		descrs["height"] = "height:交易区块高度"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "tx_nonce")
		delete(params, "amount")
		delete(params, "fee")
		delete(params, "signature")
		delete(params, "created")
		delete(params, "payload")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `查询满足条件的交易信息
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "交易查询",
			Descr:  descrStr,
			API:    "/transaction-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, transactionCtrl.Get)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "tx_nonce")
		delete(params, "created")
		delete(params, "hash")
		delete(params, "height")
		params["tx_type"] = 0
		paramsBytes, _ = json.Marshal(params)
		descrStr = `一对一转账
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "一对一转账",
			Descr:  descrStr,
			API:    "/transaction-single",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, transactionCtrl.Post)

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "tx_nonce")
		delete(params, "payload")
		delete(params, "created")
		delete(params, "hash")
		delete(params, "height")
		params["tx_type"] = 0
		paramsBytes, _ = json.Marshal(params)
		descrStr = `一对多转账
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "一对多转账",
			Descr:  descrStr,
			API:    "/transaction-multiple",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, transactionCtrl.Post)
	}

	//***********************************************************************************
	//权限管理 --- 用户账户登录与验证
	//***********************************************************************************
	{
		action := "权限管理"
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
		descrs["created"] = "created:创建时间"
		descrs["updated"] = "updated:更新时间"

		var params map[string]interface{}
		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "metadata")
		delete(params, "policy_id")
		delete(params, "status_id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ := json.Marshal(params)
		descrStr := `用户账户登录与验证
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户登录与验证",
			Descr:  descrStr,
			API:    "/user-login",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, userCtrl.Login)

		policyCtrl := controllers.NewPolicyController()

		json.Unmarshal(bytes, &params)
		delete(params, "id")
		delete(params, "name")
		delete(params, "password")
		delete(params, "metadata")
		delete(params, "status_id")
		delete(params, "created")
		delete(params, "updated")
		paramsBytes, _ = json.Marshal(params)
		descrStr = `用户账户分级分类管理
所需参数解析:
`
		for k := range params {
			descrStr += fmt.Sprintln(descrs[k])
		}
		policys = append(policys, &policy.Policy{
			Name:   "用户账户分级分类管理",
			Descr:  descrStr,
			API:    "/multiple-policy-get",
			Action: action,
			Params: string(paramsBytes),
		})
		handlers = append(handlers, policyCtrl.MultipleGet)
	}
}
