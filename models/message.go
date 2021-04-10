package models

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/validate"
)

// Message struct
type Message struct {
	ID         int64     `json:"id,omitempty"`
	TemplateID int64     `json:"templateId,omitempty"`
	ChannelID  int64     `json:"channelId,omitempty"`
	Channel    string    `json:"channel,omitempty"`
	Recipient  string    `json:"recipient,omitempty"`
	Subject    string    `json:"subject,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreatedOn  time.Time `json:"-"`
	ExpireOn   time.Time `json:"-"`
	Operation  string    `json:"operation,omitempty"`
	StatusID   int64     `json:"statusId,omitempty"`
	Status     string    `json:"status,omitempty"`
	Attempts   int64     `json:"attempts,omitempty"`
	Priority   int64     `json:"priority,omitempty"`
	Reference  string    `json:"reference,omitempty"`
	Params     []string  `json:"params,omitempty"`
}

// PrepareInput initiliazes the Message request object
func (m *Message) PrepareInput(r *http.Request) {
	// All message have expiry date of 1 week
	var expire = time.Now()
	expire.AddDate(0, 0, 7)

	m.ID = 0
	m.StatusID = 1
	m.Attempts = 0
	m.Priority = setPriority(m.ChannelID)
	m.CreatedOn = time.Now()
	m.ExpireOn = expire
}

// Validate check if request is valid
func (m *Message) Validate(action string) error {
	switch strings.ToLower(action) {
	case "new":
		if m.ChannelID == 0 {
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
		if !validate.Email(m.Recipient) {
			return errors.New("invalid email address")
		}
		return nil

	case "new-template":
		if m.TemplateID == 0 {
			return errors.New("required Template")
		}
		if m.ChannelID == 0 {
			return errors.New("required Channel")
		}
		if m.Recipient == "" {
			return errors.New("required Recipient")
		}
		if !validate.Email(m.Recipient) {
			return errors.New("invalid email address")
		}
		return nil
	default:
		return nil
	}
}

// PrepareOutput initiliazes the Message request object
func (m *Message) PrepareOutput(message Message) {
	m.ID = message.ID
}

// setPriority returns the channel priority
func setPriority(channelID int64) int64 {
	switch channelID {
	case 1:
		return 1
	case 2:
		return 2
	default:
		return 3
	}
}
