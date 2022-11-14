package models

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Message struct
type Message struct {
	ID         string    `json:"id,omitempty"`
	ChannelID  string    `json:"channelId,omitempty"`
	Recipient  string    `json:"recipient,omitempty"`
	Subject    string    `json:"subject,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreatedOn  time.Time `json:"-"`
	ExpireOn   time.Time `json:"-"`
	Operation  string    `json:"operation,omitempty"`
	Status     string    `json:"status,omitempty"`
	Attempts   int64     `json:"attempts,omitempty"`
	Priority   int64     `json:"priority,omitempty"`
	Reference  string    `json:"reference,omitempty"`
	TemplateID string    `json:"templateId,omitempty"`
	Params     []string  `json:"params,omitempty"`
}

// PrepareInput initiliazes the Message request object
func (m *Message) PrepareInput(r *http.Request) {
	// All message have expiry date of 1 week
	var expire = time.Now()
	expire.AddDate(0, 0, 7)

	m.Status = "new"
	m.Attempts = 0
	m.CreatedOn = time.Now()
	m.ExpireOn = expire
}

// Validate check if request is valid
func (m *Message) Validate(action string) error {
	switch strings.ToLower(action) {
	case "new":
		if m.ChannelID == "" {
			return errors.New("required Channel")
		}
		if m.Recipient == "" {
			return errors.New("required Recipient")
		}
		if m.Subject == "" {
			return errors.New("required Subject")
		}
		if m.Content == "" {
			return errors.New("required Content")
		}
		if m.Recipient == "" {
			return errors.New("required email address")
		}
		return nil
	default:
		return nil
	}
}

// PrepareOutput initiliazes the Message request object
func (m *Message) PrepareOutput(message Message) Message {
	res := Message{}
	res.ID = message.ID
	return res
}
