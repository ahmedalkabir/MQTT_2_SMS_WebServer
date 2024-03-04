package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mqtt_2_sms_webserver/config"
	mqtt_task "mqtt_2_sms_webserver/internal/adapter/mqtt"
	"mqtt_2_sms_webserver/internal/core/port"
	"mqtt_2_sms_webserver/internal/core/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Services struct {
	smsService port.SMSService
}

func BuildServices(sms port.SMSService) *Services {
	return &Services{
		smsService: sms,
	}
}

type API struct {
	config      *config.Settings
	logger      *zerolog.Logger
	services    *Services
	mqtt_client *mqtt_task.MQTTTask
}

func BuildAPI() *API {
	cfg, err := config.NewParsedConfig()
	if err != nil {
		panic(err)
	}

	// MQTT ========================================
	mqtt_otps := mqtt.NewClientOptions().
		AddBroker("tcp://broker.emqx.io:1883").
		SetClientID("SMSxGateway")

	mqtt_client := mqtt_task.NewMQTTTask(mqtt_otps)

	go mqtt_client.ConnectToBroker()

	// ==============================================

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.TimestampFieldName = "t"
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("app", cfg.AppCode).
		Logger()

	// services
	smsService := service.NewSMSService(cfg, mqtt_client)

	srvs := BuildServices(smsService)

	return &API{
		config:      cfg,
		logger:      &logger,
		services:    srvs,
		mqtt_client: mqtt_client,
	}
}

func (api *API) Serve(g *gin.Engine) error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", api.config.Config.Port),
		Handler:      api.Routes(g),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(api.logger, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	api.logger.Printf(fmt.Sprintf("starting server at http://localhost:%s", srv.Addr))
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		api.mqtt_client.Shutdown()

		api.logger.Printf(fmt.Sprintf("shutting down server, signal: %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	api.logger.Printf(fmt.Sprintf("stopped server, Addr: http://localhost%s", srv.Addr))
	return nil
}

func (api *API) Logger() *zerolog.Logger {
	return api.logger
}
