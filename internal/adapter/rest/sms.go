package rest

import (
	"fmt"
	"mqtt_2_sms_webserver/internal/adapter/rest/request"
	"mqtt_2_sms_webserver/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct {
	svc port.SMSService
}

func NewSMSHandler(g *gin.RouterGroup, service port.SMSService) *SMSHandler {
	handler := &SMSHandler{
		svc: service,
	}

	g.POST("/smsx", handler.SendSMS)

	return handler
}

func (h *SMSHandler) SendSMS(c *gin.Context) {
	sms_request := request.SMSSendRequest{}

	if err := c.BindJSON(&sms_request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.svc.SendSMS(c.Request.Context(), sms_request.Phone, sms_request.Message)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("the message has been sent to %s", sms_request.Phone),
	})
}
