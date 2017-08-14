package router

import (
	"net/http"

	gin "gopkg.in/gin-gonic/gin.v1"
)

var router *gin.Engine

func init() {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	router = gin.Default()
	router.Static("/js", "./public/js")
	router.Static("/css", "./public/css")
	router.Static("/assets", "./public/assets")
	router.LoadHTMLGlob("public/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

//Run start routers
func Run(addr ...string) error {
	return router.Run(addr...)
}

// IController  interface for restful api method
type IController interface {
	Get(*gin.Context)
	Post(*gin.Context)
	Put(*gin.Context)
	Delete(*gin.Context)
}

// RegisterController register restful api
func RegisterController(url string, ctrl IController) {
	RegisterGet(url, ctrl.Get)
	RegisterPost(url, ctrl.Post)
	RegisterPut(url, ctrl.Put)
	RegisterDelete(url, ctrl.Delete)
}

// RegisterGet register restful api
func RegisterGet(url string, function func(c *gin.Context)) {
	router.GET(url, function)
}

// RegisterPost register restful api
func RegisterPost(url string, function func(c *gin.Context)) {
	router.POST(url, function)
}

// RegisterPut register restful api
func RegisterPut(url string, function func(c *gin.Context)) {
	router.PUT(url, function)
}

// RegisterDelete register restful api
func RegisterDelete(url string, function func(c *gin.Context)) {
	router.DELETE(url, function)
}
