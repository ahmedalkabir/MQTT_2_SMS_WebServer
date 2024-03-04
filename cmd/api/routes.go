package api

import (
	"mqtt_2_sms_webserver/internal/adapter/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) Routes(g *gin.Engine) http.Handler {

	v1 := g.Group("/api/v1")

	rest.NewSMSHandler(v1, api.services.smsService)

	return g
}
