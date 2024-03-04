package main

import (
	"mqtt_2_sms_webserver/cmd/api"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	app := api.BuildAPI()

	if err := app.Serve(g); err != nil {
		app.Logger().Fatal().Err(err)
	}
}
