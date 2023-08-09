package mqtt

import (
	"fmt"
	"log"
	"time"

	"github.com/brestmatias/iot-libs/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	Client MQTT.Client
}

func (MqttClient) New(config *config.ConfigFile, brokerIp string) *MqttClient {
	r := MqttClient{}
	r.buildClient(config, brokerIp)
	return &r
}

func (m *MqttClient) buildClient(config *config.ConfigFile, brokerIp string) {
	//method := "buildClient"
	o := MQTT.NewClientOptions()
	o.AddBroker(fmt.Sprintf("tcp://%v:1883", brokerIp))
	o.SetClientID(config.Mqtt.ClientId)
	o.SetUsername(config.Mqtt.UserName)
	if config.Mqtt.PingTimeOut != "" {
		x, _ := time.ParseDuration(config.Mqtt.PingTimeOut)
		o.SetPingTimeout(x)
	}
	if config.Mqtt.KeepAlive != "" {
		k, _ := time.ParseDuration(config.Mqtt.KeepAlive)
		o.SetKeepAlive(k)
	}

	m.Client = MQTT.NewClient(o)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (m *MqttClient) Publish(topic string, payload interface{}) error {
	token := m.Client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MqttClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) error {
	log.Printf("Subscribing to %s", topic)
	if token := m.Client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Printf("Error subscribing to %s. error: %s", topic, token.Error())
		return token.Error()
	}
	return nil
}
