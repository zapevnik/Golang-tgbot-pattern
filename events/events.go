package events

import (
	"errors"
	"testBot/api"
	"testBot/er"
)

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Proc interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Docx
)

type Event struct {
	Type     Type
	ChatID   int
	Username string
	Text     string
}

type Processor struct {
	tg     *api.Client
	offset int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

func New(client *api.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(limit int) ([]Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, er.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event Event) error {
	switch event.Type {
	case Message:
		return p.processMessage(event)
	//case Docx:
	//	return p.processDocx(event)
	default:
		return er.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event Event) error {

	if err := p.doCmd(event); err != nil {
		return er.Wrap("can't process message", err)
	}

	return nil
}

func event(upd api.Update) Event {
	updType := fetchType(upd)

	res := Event{
		Type:     updType,
		Text:     fetchText(upd),
		ChatID:   upd.Message.Chat.ID,
		Username: upd.Message.From.Username,
	}

	return res
}

func fetchText(upd api.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd api.Update) Type {
	if upd.Message.Text != "" {
		return Message
	} else if upd.Message.Document.MimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		return Docx
	}

	return Unknown
}
