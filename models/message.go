package models

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/jwt"
	"github.com/greatfocus/gf-frame/validate"
)

// Message struct
type Message struct {
	ID        int64     `json:"id,omitempty"`
	ChannelID int64     `json:"channelId,omitempty"`
	Channel   string    `json:"channel,omitempty"`
	Recipient string    `json:"recipient,omitempty"`
	Subject   string    `json:"subject,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedBy int64     `json:"createdBy,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	ExpireOn  time.Time `json:"expireOn,omitempty"`
	StatusID  int64     `json:"statusId,omitempty"`
	Status    string    `json:"status,omitempty"`
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
	m.Priority = setPriority(m.ChannelID)

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
		if m.ChannelID == 0 {
			return errors.New("Required Channel")
		}
		if m.Recipient == "" {
			return errors.New("Required Recipient")
		}
		if m.Subject == "" {
			return errors.New("Required Subject")
		}
		if m.Content == "" {
			return errors.New("Required Content")
		}
		if !validate.Email(m.Recipient) {
			return errors.New("Invalid email address")
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
