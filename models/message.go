package models

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/jwt"
)

// Message struct
type Message struct {
	ID        int64     `json:"id,omitempty"`
	Channel   string    `json:"channel,omitempty"`
	Recipient string    `json:"recipient,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedBy int64     `json:"createdBy,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	ExpireOn  time.Time `json:"expireOn,omitempty"`
	StatusID  int64     `json:"statusId,omitempty"`
	Attempts  int64     `json:"attempts,omitempty"`
	Priority  int64     `json:"priority,omitempty"`
	RefID     int64     `json:"refId,omitempty"`
}

// PrepareInput initiliazes the Message request object
func (m *Message) PrepareInput(r *http.Request) error {
	// All message have expiry date of 1 week
	var expire = time.Now()
	expire.AddDate(0, 0, 7)

	m.ID = 0
	m.StatusID = 1
	m.Attempts = 0
	m.Priority = setPriority(m.Channel)

	m.CreatedOn = time.Now()
	m.ExpireOn = expire
	userID, err := jwt.ExtractTokenID(r)
	if err != nil {
		return errors.New("Invalid token")
	}

	m.CreatedBy = userID
	return nil
}

// Validate check if request is valid
func (m *Message) Validate(action string) error {
	switch strings.ToLower(action) {
	case "new":
		if m.Channel == "" {
			return errors.New("Required Channel")
		}
		if m.Recipient == "" {
			return errors.New("Required Recipient")
		}
		if m.Content == "" {
			return errors.New("Required Content")
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
func setPriority(channel string) int64 {
	switch strings.ToLower(channel) {
	case "sms":
		return 1
	case "email":
		return 2
	default:
		return 3
	}
}
