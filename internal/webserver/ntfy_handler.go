package webserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ntfyHandler function moved to this file
func ntfyHandler(messageChannel chan Message) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(rw, "Invalid topic in URL", http.StatusBadRequest)
			return
		}
		topic := pathParts[2]

		title := r.Header.Get("Title")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Error reading body", http.StatusBadRequest)
			return
		}

		log.Printf("Received notification for topic '%s': Title: '%s', Message: '%s'", topic, title, string(body))
		rw.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(rw, "Notification for topic '%s' received", topic)

		messageChannel <- Message{Topic: topic, Body: string(body), Title: title}
	}
	return http.HandlerFunc(fn)
}
