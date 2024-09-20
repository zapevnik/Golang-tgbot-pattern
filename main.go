package main

import (
	"flag"
	"log"
	"testBot/api"
	"testBot/consumer"
	"testBot/events"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	eventsProcessor := events.New(
		api.New(tgBotHost, mustToken()),
	)

	log.Print("service started")

	con := consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := con.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
