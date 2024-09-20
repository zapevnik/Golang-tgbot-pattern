package events

import (
	"log"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(event Event) error {
	event.Text = strings.TrimSpace(event.Text)

	log.Printf("got new command '%s' from '%s", event.Text, event.Username)

	switch event.Text {
	case HelpCmd:
		return p.sendHelp(event.ChatID)
	case StartCmd:
		return p.sendHello(event.ChatID)
	default:
		return p.tg.SendMessage(event.ChatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
