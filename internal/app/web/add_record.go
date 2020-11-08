package web

import (
	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func AddDDNSRecord(ctx *gin.Context) {
	domain, ok := ctx.GetPostForm("domain")
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("undefined domain"))
		return
	}

	updatePeriod, ok := ctx.GetPostForm("updatePeriod")
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("undefined updatePeriod"))
	}

	period, err := time.ParseDuration(updatePeriod)
	if err != nil {
		err = errors.Wrap(err, "invalid updatePeriod format")
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	event := bus.AddDomainRecord{
		Domain:       domain,
		UpdatePeriod: period,
	}
	bus.Get().Publish(event)
	ctx.Redirect(http.StatusFound, "/")
}
