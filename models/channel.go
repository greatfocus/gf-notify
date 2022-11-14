package models

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Channel struct
type Channel struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Key       string    `json:"key,omitempty"`
	URL       string    `json:"url,omitempty"`
	User      string    `json:"user,omitempty"`
	Pass      string    `json:"pass,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	UpdatedOn time.Time `json:"updatedOn,omitempty"`
	Enabled   bool      `json:"enabled,omitempty"`
}

// PrepareChannel initiliazes the channel request object
func (c *Channel) PrepareChannel(r *http.Request) {
	c.CreatedOn = time.Now()
	c.UpdatedOn = time.Now()
}

// ValidateChannel check if request is valid
func (c *Channel) ValidateChannel(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if c.Name == "" {
			return errors.New("name is required")
		}
		if c.Key == "" {
			return errors.New("key is required")
		}
		if c.URL == "" {
			return errors.New("url is required")
		}
		if c.User == "" {
			return errors.New("user is required")
		}
		if c.Pass == "" {
			return errors.New("pass is required")
		}
		return nil
	case "update":
		if c.ID == "" {
			return errors.New("id is required")
		}
		if c.Name == "" {
			return errors.New("name is required")
		}
		if c.Key == "" {
			return errors.New("key is required")
		}
		if c.URL == "" {
			return errors.New("url is required")
		}
		if c.User == "" {
			return errors.New("user is required")
		}
		if c.Pass == "" {
			return errors.New("pass is required")
		}
		return nil
	default:
		return nil
	}
}
