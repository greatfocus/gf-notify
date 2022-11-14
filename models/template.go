package models

import (
	"errors"
	"strings"
	"time"
)

// Template struct
type Template struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Key       string    `json:"key,omitempty"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body,omitempty"`
	CreatedOn time.Time `json:"-"`
	UpdatedOn time.Time `json:"-"`
	Enabled   bool      `json:"-"`
	Deleted   bool      `json:"-"`
}

// PrepareTempate initiliazes the Template request object
func (t *Template) PrepareTempate() {
	t.UpdatedOn = time.Now()
	t.CreatedOn = time.Now()
}

// ValidateTemplate check if request is valid
func (t *Template) ValidateTemplate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if t.ID == "" {
			return errors.New("required ID")
		}
		if t.Name == "" {
			return errors.New("required Name")
		}
		if t.Key == "" {
			return errors.New("required Key")
		}
		if t.Subject == "" {
			return errors.New("required Subject")
		}
		if t.Body == "" {
			return errors.New("required Body")
		}
		return nil

	case "create":
		if t.Name == "" {
			return errors.New("required Name")
		}
		if t.Key == "" {
			return errors.New("required Key")
		}
		if t.Subject == "" {
			return errors.New("required Subject")
		}
		if t.Body == "" {
			return errors.New("required Body")
		}
		return nil
	default:
		return errors.New("invalid validation operation")
	}
}

// PrepareTemplateOutput prepare the template to output
func (t *Template) PrepareTemplateOutput(temp Template) Template {
	res := Template{}
	t.ID = temp.ID
	t.Key = temp.Key
	return res
}
