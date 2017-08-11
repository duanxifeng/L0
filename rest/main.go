package main

import (
	"github.com/bocheninc/L0/rest/api"
	_ "github.com/bocheninc/L0/rest/model"
)

func main() {
	api.Run(":8000")
}
