package main

import (
	"log"
	"os"
	"strconv"
)

func getEnvInt(key string) int {
	var variable string = os.Getenv(key)
	if len(variable) == 0 {
		variable = "0"
	}
	val, err := strconv.Atoi(variable)
	if err != nil {
		log.Fatal("Failed to convert env var to integer")
	}
	return val
}

func logEnvs() {
	log.Printf(`Environment variables
	- PORT: %v`,
		PORT)
}
