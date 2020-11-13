package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/exepirit/cf-ddns/internal/bus"
)

func AddDnsBinding(ctx *gin.Context) {
	domain, ok := ctx.GetPostForm("domain")
	if !ok {
		ctx.JSON(http.StatusBadRequest, "required field \"domain\"")
		return
	}

	updatePeriod, ok := ctx.GetPostForm("updatePeriod")
	if !ok {
		ctx.JSON(http.StatusBadRequest, "required field \"updatePeriod\"")
		return
	}

	d, err := time.ParseDuration(updatePeriod)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid time.Period format")
		return
	}

	bus.Get().Publish(bus.AddDomainBinding{
		Domain:       domain,
		UpdatePeriod: d,
	})
	ctx.Redirect(http.StatusFound, "/")
}
