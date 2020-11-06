package web

import (
	"github.com/exepirit/cf-ddns/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecordsPage(ctx *gin.Context) {
	repo := repository.Get()
	records := repo.GetAll()

	ctx.HTML(http.StatusOK, "records.html", gin.H{
		"Records": records,
	})
}
