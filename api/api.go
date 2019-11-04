package api

import (
	"github.com/gin-gonic/gin"
)

type Api struct {
	engine *gin.Engine
}

func New(debug bool) *Api {
	api := &Api{}
	api.init(debug)
	return api
}

func (a *Api) init(debug bool) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()
	a.engine = gin.Default()
	a.engine.GET("/ping", handlePing)
}

func (a *Api) Start() {
	a.engine.Run()
}
