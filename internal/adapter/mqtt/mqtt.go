package mqtt_task

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTTask struct {
	mqtt_opts   *mqtt.ClientOptions
	mqtt_client mqtt.Client
	shutdown    chan bool
}

func NewMQTTTask(opts *mqtt.ClientOptions) *MQTTTask {
	return &MQTTTask{
		mqtt_opts: opts,
		shutdown:  make(chan bool),
	}
}

func (t *MQTTTask) ConnectToBroker() {
	fmt.Println("[ConnectToBroker]")
	c := mqtt.NewClient(t.mqtt_opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	t.mqtt_client = c
	status := <-t.shutdown

	if status {
		fmt.Println("[shutdown connection]")
		c.Disconnect(10)
	}
}

func (t *MQTTTask) PublishMessage(topic string, message interface{}) {
	go func() {
		token := t.mqtt_client.Publish(topic, 0, false, message)
		token.Wait()
	}()
}

func (t *MQTTTask) Shutdown() {
	t.shutdown <- true
}
