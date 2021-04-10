package models

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Channel struct
type Channel struct {
	ID         int64     `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	StaticName string    `json:"staticName,omitempty"`
	Priority   int64     `json:"priority,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
	UpdatedOn  time.Time `json:"updatedOn,omitempty"`
	Enabled    bool      `json:"enabled,omitempty"`
}

// PrepareChannel initiliazes the channel request object
func (c *Channel) PrepareChannel(r *http.Request) {
	c.CreatedOn = time.Now()
	c.UpdatedOn = time.Now()
}

// ValidateChannel check if request is valid
func (c *Channel) ValidateChannel(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.ID == 0 {
			return errors.New("required ID")
		}
		if c.Priority == 0 {
			return errors.New("required Priority")
		}
		return nil
	default:
		return nil
	}
}
