package models

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/jwt"
)

// Template struct
type Template struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	StaticName  string    `json:"staticName,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	Body        string    `json:"body,omitempty"`
	ParamsCount int64     `json:"paramsCount,omitempty"`
	CreatedBy   int64     `json:"-"`
	CreatedOn   time.Time `json:"-"`
	UpdatedBy   int64     `json:"-"`
	UpdatedOn   time.Time `json:"-"`
	Enabled     bool      `json:"-"`
	Deleted     bool      `json:"-"`
}

// PrepareTempate initiliazes the Template request object
func (c *Template) PrepareTempate(r *http.Request) error {
	c.UpdatedOn = time.Now()
	c.CreatedOn = time.Now()

	userID, err := jwt.ExtractTokenID(r)
	if err != nil {
		return errors.New("Invalid token")
	}

	c.CreatedBy = userID
	c.UpdatedBy = userID
	return nil
}

// ValidateTemplate check if request is valid
func (c *Template) ValidateTemplate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if c.ID == 0 {
			return errors.New("Required ID")
		}
		if c.Name == "" {
			return errors.New("Required Name")
		}
		if c.StaticName == "" {
			return errors.New("Required StaticName")
		}
		if c.Subject == "" {
			return errors.New("Required Subject")
		}
		if c.Body == "" {
			return errors.New("Required Body")
		}
		return nil

	case "add":
		if c.Name == "" {
			return errors.New("Required Name")
		}
		if c.StaticName == "" {
			return errors.New("Required StaticName")
		}
		if c.Subject == "" {
			return errors.New("Required Subject")
		}
		if c.Body == "" {
			return errors.New("Required Body")
		}
		return nil
	default:
		return errors.New("Invalid validation operation")
	}
}
