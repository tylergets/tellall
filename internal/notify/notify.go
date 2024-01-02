package notify

import (
	"github.com/containrrr/shoutrrr"
	"github.com/gen2brain/beeep"
	"github.com/tylergets/tellall/internal/config"
	"github.com/tylergets/tellall/internal/mqtt"
	"log"
	"strings"
)

type Notifier struct {
	config *config.Config
	client *mqtt.Client
}

const (
	HTTP = "http"
	MQTT = "mqtt"
)

type Notification struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Topic    string `json:"topic"`
	SentFrom string `json:"sentFrom"`
	Source   string `json:"source"`
}

func NewNotifier(config *config.Config, client *mqtt.Client) *Notifier {
	return &Notifier{
		client: client,
		config: config,
	}
}

func (n *Notifier) SendNotificationToAll(notification Notification) {
	notification.Topic = "all"
	n.SendNotificationToTopic(notification)
}

func (n *Notifier) SendNotificationToTopic(notification Notification) {

	if notification.Source == HTTP {
		n.client.PublishMessage(notification.Topic, notification.Title, notification.Body)
	}

	for _, listener := range n.config.Listeners {
		// if it starts with "notify-send://" then use notify-send, otherwise pass off to shoutarr
		if strings.HasPrefix(listener, "notify-send://") {
			err := beeep.Notify(notification.Title, notification.Body, "assets/information.png")
			if err != nil {
				panic(err)
			}
		} else {
			log.Println("Sending notification to " + listener)
			err := shoutrrr.Send(listener+"?title="+notification.Title, notification.Body)
			if err != nil {
				panic(err)
			}
		}
	}

}
