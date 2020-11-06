package web

import "github.com/gin-gonic/gin"

// New creates new gin.Engine and attach routes.
func New() *gin.Engine {
	engine := gin.New()
	engine.Static("/css", "web/css")
	engine.LoadHTMLGlob("web/templates/*")
	engine.Use(gin.Recovery(), gin.Logger())

	engine.GET("/", RecordsPage)

	return engine
}
