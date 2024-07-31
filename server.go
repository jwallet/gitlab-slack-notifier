package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handshake struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

func setupServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
	})

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("OK"))
		fmt.Println("healthcheck")
	})

	mux.HandleFunc("/gitlab-webhook", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		receivedSignature := r.Header.Get("X-Gitlab-Token")
		if receivedSignature != GITLAB_WEBHOOK_SECRET_TOKEN {
			fmt.Printf("Invalid secret token, received ''%v'', expected ''%v''\n", receivedSignature, GITLAB_WEBHOOK_SECRET_TOKEN)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var webhookEvent GitLabWebhookEvent
		json.Unmarshal(body, &webhookEvent)

		err = handleGitLabWebhook(webhookEvent)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNoContent)
		}

		fmt.Println("---------------")
	})

	mux.HandleFunc("/slack-events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payload SlackPayload
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(body, &payload)

		if payload.Type == "event_callback" {
			err = handleSlackEvent(payload.Event)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNoContent)
			}
		} else if payload.Type == "url_verification" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(payload.Challenge))
		}
		fmt.Println("---------------")
	})

	fmt.Println("Starting server...")

	// Determine port for HTTP service.
	port := PORT
	if port == 0 {
		port = 8080
		fmt.Printf("Defaulting to port %v\n", port)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	server.SetKeepAlivesEnabled(false)

	fmt.Printf("Server listening on localhost:%v\n", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
