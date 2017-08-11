package api

import (
	"encoding/json"
	"time"

	gin "gopkg.in/gin-gonic/gin.v1"

	"github.com/bocheninc/L0/rest/controllers"
	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/account"
	"github.com/bocheninc/L0/rest/model/table/history"
	"github.com/bocheninc/L0/rest/model/table/policy"
	"github.com/bocheninc/L0/rest/model/table/status"
	"github.com/bocheninc/L0/rest/model/table/user"
	"github.com/bocheninc/L0/rest/router"
)

//Run start api service
func Run(addr ...string) {
	users := make([]*user.User, 0)
	users = append(users, &user.User{
		Name:     "admin",
		PassWord: "admin",
	})

	//init status table data
	statuses := make([]*status.Status, 0)
	statuses = append(statuses, &status.Status{
		Name:    "in",
		Descr:   "pend",
		Created: time.Now(),
		Updated: time.Now(),
	})
	statuses = append(statuses, &status.Status{
		Name:    "on",
		Descr:   "normal",
		Created: time.Now(),
		Updated: time.Now(),
	})
	statuses = append(statuses, &status.Status{
		Name:    "off",
		Descr:   "expired",
		Created: time.Now(),
		Updated: time.Now(),
	})

	//init policy table data
	policys := make([]*policy.Policy, 0)
	handlers := make([]func(c *gin.Context), 0)

	//***********************************************************************************
	//*************************************history***************************************
	//***********************************************************************************
	{
		history := history.NewHistory()
		historyCtrl := controllers.NewHistoryController()
		policys = append(policys, &policy.Policy{
			Name:   "history-get",
			Descr:  history.GetDescr(),
			API:    "/history-get",
			Action: history.TableName(),
			Params: history.GetParams(),
		})
		handlers = append(handlers, historyCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:   "history-post",
			Descr:  history.PostDescr(),
			API:    "/history-post",
			Action: history.TableName(),
			Params: history.PostParams(),
		})
		handlers = append(handlers, historyCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:   "history-put",
			Descr:  history.PutDescr(),
			API:    "/history-put",
			Action: history.TableName(),
			Params: history.PutParams(),
		})
		handlers = append(handlers, historyCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:   "history-delete",
			Descr:  history.DeleteDescr(),
			API:    "/history-delete",
			Action: history.TableName(),
			Params: history.DeleteParams(),
		})
		handlers = append(handlers, historyCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************status******************************************
	//***********************************************************************************
	{
		status := status.NewStatus()
		statusCtrl := controllers.NewStatusController()
		policys = append(policys, &policy.Policy{
			Name:   "status-get",
			Descr:  status.GetDescr(),
			API:    "/statuse-get",
			Action: status.TableName(),
			Params: status.GetParams(),
		})
		handlers = append(handlers, statusCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:   "status-post",
			Descr:  status.PostDescr(),
			API:    "/status-post",
			Action: status.TableName(),
			Params: status.PostParams(),
		})
		handlers = append(handlers, statusCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:   "status-put",
			Descr:  status.PutDescr(),
			API:    "/status-put",
			Action: status.TableName(),
			Params: status.PutParams(),
		})
		handlers = append(handlers, statusCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:   "status-delete",
			Descr:  status.DeleteDescr(),
			API:    "/status-delete",
			Action: status.TableName(),
			Params: status.DeleteParams(),
		})
		handlers = append(handlers, statusCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************policy******************************************
	//***********************************************************************************
	{
		opolicy := policy.NewPolicy()
		policyCtrl := controllers.NewPolicyController()
		policys = append(policys, &policy.Policy{
			Name:   "policy-get",
			Descr:  opolicy.GetDescr(),
			API:    "/policy-get",
			Action: opolicy.TableName(),
			Params: opolicy.GetParams(),
		})
		handlers = append(handlers, policyCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:   "policy-post",
			Descr:  opolicy.PostDescr(),
			API:    "/policy-post",
			Action: opolicy.TableName(),
			Params: opolicy.PostParams(),
		})
		handlers = append(handlers, policyCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:   "policy-put",
			Descr:  opolicy.PutDescr(),
			API:    "/policy-put",
			Action: opolicy.TableName(),
			Params: opolicy.PutParams(),
		})
		handlers = append(handlers, policyCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:   "policy-delete",
			Descr:  opolicy.DeleteDescr(),
			API:    "/policy-delete",
			Action: opolicy.TableName(),
			Params: opolicy.DeleteParams(),
		})
		handlers = append(handlers, policyCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************user******************************************
	//***********************************************************************************
	{
		user := user.NewUser()
		userCtrl := controllers.NewUserController()
		policys = append(policys, &policy.Policy{
			Name:   "user-get",
			Descr:  user.GetDescr(),
			API:    "/user-get",
			Action: user.TableName(),
			Params: user.GetParams(),
		})
		handlers = append(handlers, userCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:   "user-post",
			Descr:  user.PostDescr(),
			API:    "/user-post",
			Action: user.TableName(),
			Params: user.PostParams(),
		})
		handlers = append(handlers, userCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:   "user-put",
			Descr:  user.PutDescr(),
			API:    "/user-put",
			Action: user.TableName(),
			Params: user.PutParams(),
		})
		handlers = append(handlers, userCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:   "user-delete",
			Descr:  user.DeleteDescr(),
			API:    "/user-delete",
			Action: user.TableName(),
			Params: user.DeleteParams(),
		})
		handlers = append(handlers, userCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************account***************************************
	//***********************************************************************************
	{
		account := account.NewAccount()
		accountCtrl := controllers.NewAccountController()
		policys = append(policys, &policy.Policy{
			Name:   "account-get",
			Descr:  account.GetDescr(),
			API:    "/account-get",
			Action: account.TableName(),
			Params: account.GetParams(),
		})
		handlers = append(handlers, accountCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:   "account-post",
			Descr:  account.PostDescr(),
			API:    "/account-post",
			Action: account.TableName(),
			Params: account.PostParams(),
		})
		handlers = append(handlers, accountCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:   "account-put",
			Descr:  account.PutDescr(),
			API:    "/account-put",
			Action: account.TableName(),
			Params: account.PutParams(),
		})
		handlers = append(handlers, accountCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:   "account-delete",
			Descr:  account.DeleteDescr(),
			API:    "/account-delete",
			Action: account.TableName(),
			Params: account.DeleteParams(),
		})
		handlers = append(handlers, accountCtrl.Delete)
	}

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

	for _, user := range users {
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

	for index, opolicy := range policys {
		function := handlers[index]
		var params map[string]string
		json.Unmarshal([]byte(opolicy.Params), &params)
		action := opolicy.Action
		router.RegisterPost(opolicy.API, func(c *gin.Context) {
			for k := range params {
				params[k] = c.PostForm(k)
			}
			paramsBytes, _ := json.Marshal(params)
			// add history
			history := &history.History{
				API:     c.Request.URL.String(),
				Action:  action,
				Params:  string(paramsBytes),
				Created: time.Now(),
				Updated: time.Now(),
			}
			addHistory(history)

			function(c)
		})
	}
	router.Run(addr...)
}

func addHistory(history *history.History) {
	tx, _ := model.DB.Begin()
	if err := history.Insert(tx); err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
