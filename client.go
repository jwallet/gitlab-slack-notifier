package main

import (
	"net/http"
	"time"
)

func getClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true

	return &http.Client{
		Transport: t,
		Timeout:   10 * time.Second,
	}
}
