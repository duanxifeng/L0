package api

import (
	"encoding/json"
	"time"

	gin "gopkg.in/gin-gonic/gin.v1"

	"github.com/bocheninc/L0/rest/controllers"
	"github.com/bocheninc/L0/rest/model"
	"github.com/bocheninc/L0/rest/model/table/history"
	"github.com/bocheninc/L0/rest/model/table/policy"
	"github.com/bocheninc/L0/rest/model/table/status"
	"github.com/bocheninc/L0/rest/model/table/user"
	"github.com/bocheninc/L0/rest/router"
)

//Run start api service
func Run(addr ...string) {
	admin := &user.User{
		Name:     "admin",
		PassWord: "admin",
	}
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
		historyCtrl := controllers.NewHistoryController()
		policys = append(policys, &policy.Policy{
			Name:    "history-get",
			Descr:   "to do",
			API:     "/historys",
			Action:  "get",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, historyCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:    "history-post",
			Descr:   "to do",
			API:     "/historys",
			Action:  "post",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, historyCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:    "history-put",
			Descr:   "to do",
			API:     "/historys",
			Action:  "put",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, historyCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:    "history-delete",
			Descr:   "to do",
			API:     "/historys",
			Action:  "delete",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, historyCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************status******************************************
	//***********************************************************************************
	{
		statusCtrl := controllers.NewStatusController()
		policys = append(policys, &policy.Policy{
			Name:    "status-get",
			Descr:   "to do",
			API:     "/statuses",
			Action:  "get",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, statusCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:    "status-post",
			Descr:   "to do",
			API:     "/statuses",
			Action:  "post",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, statusCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:    "status-put",
			Descr:   "to do",
			API:     "/statuses",
			Action:  "put",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, statusCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:    "status-delete",
			Descr:   "to do",
			API:     "/statuses",
			Action:  "delete",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, statusCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************policy******************************************
	//***********************************************************************************
	{
		policyCtrl := controllers.NewPolicyController()
		policys = append(policys, &policy.Policy{
			Name:    "policy-get",
			Descr:   "to do",
			API:     "/policys",
			Action:  "get",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, policyCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:    "policy-post",
			Descr:   "to do",
			API:     "/policys",
			Action:  "post",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, policyCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:    "policy-put",
			Descr:   "to do",
			API:     "/policys",
			Action:  "put",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, policyCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:    "policy-delete",
			Descr:   "to do",
			API:     "/policys",
			Action:  "delete",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, policyCtrl.Delete)
	}

	//***********************************************************************************
	//*************************************user******************************************
	//***********************************************************************************
	{
		userCtrl := controllers.NewUserController()
		policys = append(policys, &policy.Policy{
			Name:    "user-get",
			Descr:   "to do",
			API:     "/users",
			Action:  "get",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, userCtrl.Get)

		policys = append(policys, &policy.Policy{
			Name:    "user-post",
			Descr:   "to do",
			API:     "/users",
			Action:  "post",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, userCtrl.Post)

		policys = append(policys, &policy.Policy{
			Name:    "user-put",
			Descr:   "to do",
			API:     "/users",
			Action:  "put",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, userCtrl.Put)

		policys = append(policys, &policy.Policy{
			Name:    "user-delete",
			Descr:   "to do",
			API:     "/users",
			Action:  "delete",
			Params:  "",
			Created: time.Now(),
			Updated: time.Now(),
		})
		handlers = append(handlers, userCtrl.Delete)
	}

	tx, _ := model.DB.Begin()
	for _, status := range statuses {
		if err := status.InsertOrUpdate(tx); err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	for _, opolicy := range policys {
		if err := opolicy.InsertOrUpdate(tx); err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	if err := admin.InsertOrUpdate(tx); err != nil {
		panic(err)
	}
	tx.Commit()

	for index, opolicy := range policys {
		function := handlers[index]
		var params map[string]string
		json.Unmarshal([]byte(opolicy.Params), &params)
		switch action := opolicy.Action; action {
		case "get":
			router.RegisterGet(opolicy.API, func(c *gin.Context) {
				// for k := range params {
				// 	params[k] = c.PostForm(k)
				// }
				// paramsBytes, _ := json.Marshal(params)
				// // add history
				// history := &history.History{
				// 	API:     c.Request.URL.String(),
				// 	Action:  action,
				// 	Params:  string(paramsBytes),
				// 	Created: time.Now(),
				// 	Updated: time.Now(),
				// }
				// addHistory(history)

				function(c)
			})
		case "post":
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
		case "put":
			router.RegisterPut(opolicy.API, func(c *gin.Context) {
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
		case "delete":
			router.RegisterDelete(opolicy.API, func(c *gin.Context) {
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
		default:
		}
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
