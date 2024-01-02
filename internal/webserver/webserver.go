package webserver

import (
	"github.com/tylergets/tellall/internal/config"
	"log"
	"net/http"
)

type Message struct {
	Topic string
	Body  string
	Title string
}

func SetupServerAndListen(config *config.Config) chan Message {

	messageChannel := make(chan Message, 100)

	if config.HttpServer.Enabled {
		http.Handle("/ntfy/", ntfyHandler(messageChannel))
		http.Handle("/gotify/", gotifyHandler(messageChannel))

		go func() {
			log.Println("Starting HTTP server on http://" + config.HttpServer.Host + ":" + config.HttpServer.Port)

			err := http.ListenAndServe(config.HttpServer.Host+":"+config.HttpServer.Port, nil)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()

		log.Println("Listening on port " + config.HttpServer.Port)

	}

	return messageChannel
}
