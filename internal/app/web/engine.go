package web

import (
	"github.com/exepirit/cf-ddns/internal/repository"
	"github.com/gin-gonic/gin"
)

var recordsRepo repository.RecordGetter

// New creates new gin.Engine and attach routes.
func New(repo repository.RecordGetter) *gin.Engine {
	recordsRepo = repo

	engine := gin.New()
	engine.Static("/css", "web/css")
	engine.LoadHTMLGlob("web/templates/*")
	engine.Use(gin.Recovery(), gin.Logger())

	engine.GET("/", RecordsPage)

	return engine
}
