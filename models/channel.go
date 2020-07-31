package models

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/jwt"
)

// Channel struct
type Channel struct {
	ID         int64     `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	StaticName string    `json:"staticName,omitempty"`
	Priority   int64     `json:"priority,omitempty"`
	CreatedBy  int64     `json:"createdBy,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
	UpdatedBy  int64     `json:"updatedBy,omitempty"`
	UpdatedOn  time.Time `json:"updatedOn,omitempty"`
	Enabled    bool      `json:"enabled,omitempty"`
}

// PrepareChannel initiliazes the channel request object
func (c *Channel) PrepareChannel(r *http.Request) error {
	c.UpdatedOn = time.Now()

	userID, _, err := jwt.ExtractTokenID(r)
	if err != nil {
		return errors.New("Invalid token")
	}

	c.UpdatedBy = userID
	return nil
}

// ValidateChannel check if request is valid
func (c *Channel) ValidateChannel(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.ID == 0 {
			return errors.New("Required ID")
		}
		if c.Priority == 0 {
			return errors.New("Required Priority")
		}
		if c.UpdatedBy == 0 {
			return errors.New("Required UpdatedBy")
		}
		return nil
	default:
		return nil
	}
}
