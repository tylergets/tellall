package webserver

import (
	"fmt"
	"log"
	"net/http"
)

func gotifyHandler(messageChannel chan Message) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("Error parsing multipart form data: %v", err)
			http.Error(rw, "Error parsing form data", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		message := r.FormValue("message")

		if title == "" || message == "" {
			log.Printf("Empty title or message, not processing")
			http.Error(rw, "Empty title or message", http.StatusBadRequest)
			return
		}

		log.Printf("Received Gotify message: Title: '%s', Message: '%s'", title, message)
		rw.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(rw, "Gotify message received")

		messageChannel <- Message{Topic: "all", Body: message, Title: title}
	}
	return http.HandlerFunc(fn)
}
