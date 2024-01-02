package mqtt

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/enescakir/emoji"
	"github.com/tylergets/tellall/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println(emoji.SmilingFace.String() + "Connected to MQTT broker")
	}
	connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Printf("Connect lost: %v", err)
	}
)

type MessagePayload struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Topic    string `json:"topic"`
	SentFrom string `json:"sentFrom"`
}

type MqttMessage struct {
	Topic   string
	Message string
}

type Client struct {
	mqttClient mqtt.Client
	verbose    bool
	config     *config.Config
	Messages   chan MqttMessage
}

func NewClient(config *config.Config) *Client {
	if config.Debug {
		log.SetOutput(os.Stdout)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttConnection)
	opts.SetClientID("tellall")

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	msgChannel := make(chan MqttMessage)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf(emoji.Egg.String()+"Received MQTT message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		msgChannel <- MqttMessage{Topic: msg.Topic(), Message: string(msg.Payload())}
	})

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	client := &Client{
		mqttClient: mqttClient,
		verbose:    config.Debug,
		config:     config,
		Messages:   msgChannel,
	}

	return client
}

func (c *Client) WaitForMessages() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	log.Println("Waiting for messages. Press Ctrl+C to exit.")
	select {
	case <-sigc:
		log.Println("Received shutdown signal")
	}
}

func (c *Client) Subscribe(topic string, qos byte) {
	if token := c.mqttClient.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}
	log.Printf(emoji.Ear.String()+" Subscribed to %s\n", topic)
}

func (c *Client) Publish(topic string, qos byte, retained bool, message string) {
	token := c.mqttClient.Publish(topic, qos, retained, message)
	token.Wait()
	log.Printf(emoji.Megaphone.String()+" Published MQTT message: %s to topic: %s\n", message, topic)
}

func (c *Client) PublishMessage(topic string, title string, message string) {
	pubTopic := c.config.Prefix + "/device/" + topic + "/messages"

	msg := MessagePayload{}
	msg.Title = title
	msg.Body = message
	msg.Topic = topic
	msg.SentFrom = c.config.Name

	// json encode
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		return
	}

	c.Publish(pubTopic, 0, false, string(msgBytes))
}

func (c *Client) Disconnect() {
	c.mqttClient.Disconnect(250)
	log.Println("Disconnected")
}
