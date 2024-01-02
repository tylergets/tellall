package main

import (
	"encoding/json"
	"fmt"
	"github.com/tylergets/tellall/internal/config"
	"github.com/tylergets/tellall/internal/mqtt"
	"github.com/tylergets/tellall/internal/notify"
	"github.com/tylergets/tellall/internal/webserver"
	"log"
)

func main() {
	appVersion := "0.0.1"
	fmt.Println("Starting TellAll version " + appVersion)

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// print count of config.Listeners array
	log.Printf("Loaded %d listeners", len(config.Listeners))

	client := mqtt.NewClient(config)

	notifier := notify.NewNotifier(config, client)

	log.Println("Starting MQTT client")

	defer client.Disconnect()

	// Subscribe to topic.
	client.Subscribe(config.Prefix+"/device/all/messages", 0)
	client.Subscribe(config.Prefix+"/device/"+config.Name+"/messages", 0)

	http := webserver.SetupServerAndListen(config)

	go func() {
		for {
			select {

			// MQTT
			case message := <-client.Messages:

				notification := notify.Notification{}
				err := json.Unmarshal([]byte(message.Message), &notification)
				notification.Source = "mqtt"
				if err != nil {
					log.Printf("Error decoding JSON: %v", err)
				}

				if notification.SentFrom == config.Name {
					log.Printf("Ignoring message from self")
					continue
				}

				notifier.SendNotificationToTopic(notification)

			// HTTP
			case msg := <-http:

				notifier.SendNotificationToTopic(notify.Notification{
					Title:  msg.Title,
					Body:   msg.Body,
					Topic:  msg.Topic,
					Source: "http",
				})
			}
		}
	}()

	client.WaitForMessages()
}
