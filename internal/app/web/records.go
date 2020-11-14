package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecordsPage(ctx *gin.Context) {
	records := recordsRepo.GetAll()

	ctx.Header("Access-Control-Allow-Origin", "code.jquery.com, cdn.jsdelivr.net")
	ctx.HTML(http.StatusOK, "records.html", gin.H{
		"Records": records,
	})
}
