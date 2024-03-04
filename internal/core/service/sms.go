package service

import (
	"context"
	"encoding/json"
	"mqtt_2_sms_webserver/config"
	mqtt_task "mqtt_2_sms_webserver/internal/adapter/mqtt"
)

type SMSService struct {
	config      *config.Settings
	mqtt_client *mqtt_task.MQTTTask
}

func NewSMSService(config *config.Settings, mqtt_client *mqtt_task.MQTTTask) *SMSService {
	return &SMSService{
		config:      config,
		mqtt_client: mqtt_client,
	}
}

func (s *SMSService) SendSMS(ctx context.Context, phone_number, message string) error {
	topic := ""
	payload := map[string]string{
		"phone":   phone_number,
		"message": message,
	}
	byte_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	s.mqtt_client.PublishMessage(topic, byte_payload)

	return nil
}
