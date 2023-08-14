package mqtt

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"log"
	"time"
)

type MqttPublisher struct {
	SentCommands []CommandHash
	MqttClient   *MqttClient
	minInterval  time.Duration
}

type CommandHash struct {
	Topic    string
	LastHash string
	LastSent time.Time
}

func (MqttPublisher) New(mqttClient *MqttClient, minInterval time.Duration) *MqttPublisher {
	service := MqttPublisher{
		MqttClient:  mqttClient,
		minInterval: minInterval,
	}
	return &service
}

/*
Envio de comando manteniendo espacio definido por configuración
Evita inundar la cola del tópico con envío recurrente del mismo mensaje
Si el mensaje a enviar es igual al anterior, deberá cumplirse el intervalo de espacio definido
*/
func (m *MqttPublisher) SpacedPublishCommand(topic string, message interface{}) bool {
	method := "SpacedPublishCommand"

	if !m.shouldSend(topic, message) {
		return false
	}

	messageJSON, _ := json.Marshal(message)
	m.MqttClient.Publish(topic, messageJSON)
	log.Printf("[method:%v][topic:%v] Command Published", method, topic)

	return true
}

func (m *MqttPublisher) shouldSend(topic string, message interface{}) bool {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", message)))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	if len(m.SentCommands) == 0 {
		m.SentCommands = append(m.SentCommands, CommandHash{
			Topic:    topic,
			LastHash: hash,
			LastSent: time.Now(),
		})
		return true
	}

	for i, command := range m.SentCommands {
		if command.Topic == topic {
			if command.LastHash == hash {
				if diff := time.Now().Sub(command.LastSent); diff >= m.minInterval {
					(&m.SentCommands[i]).LastSent = time.Now()
					return true
				} else {
					return false
				}
			} else {
				(&m.SentCommands[i]).LastHash = hash
				(&m.SentCommands[i]).LastSent = time.Now()
				return true
			}
		}
	}
	m.SentCommands = append(m.SentCommands, CommandHash{
		Topic:    topic,
		LastHash: hash,
		LastSent: time.Now(),
	})
	return true
}

func (m *MqttPublisher) PublishCommand(topic string, message interface{}) bool {
	method := "PublishCommand"

	messageJSON, _ := json.Marshal(message)
	/*if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[method:%v] %v", method, token.Error().Error())
		return false
	}*/
	err := m.MqttClient.Publish(topic, messageJSON)
	if err == nil {
		log.Printf("[method:%v][topic:%v] error publishing command", method, topic)
		return false
	}
	log.Printf("[method:%v][topic:%v] Command Published", method, topic)

	// TODO !!!! OJO con esto que sigue, porque se desconectaba el cliente y se desuscribía
	// revisar si es necesario recheckear la conexión cada x tiempo......
	// sobre todo por las suscripciones
	//m.Client.Disconnect(250) <-- no descomentar, solo lo dejo para acordarme de investigar

	return true
}
