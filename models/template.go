package models

import (
	"errors"
	"strings"
	"time"
)

// Template struct
type Template struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	StaticName  string    `json:"staticName,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	Body        string    `json:"body,omitempty"`
	ParamsCount int64     `json:"paramsCount,omitempty"`
	CreatedOn   time.Time `json:"-"`
	UpdatedOn   time.Time `json:"-"`
	Enabled     bool      `json:"-"`
	Deleted     bool      `json:"-"`
}

// PrepareTempate initiliazes the Template request object
func (t *Template) PrepareTempate() {
	t.UpdatedOn = time.Now()
	t.CreatedOn = time.Now()
}

// ValidateTemplate check if request is valid
func (t *Template) ValidateTemplate(action string) error {
	switch strings.ToLower(action) {
	case "edit":
		if t.ID == 0 {
			return errors.New("Required ID")
		}
		if t.Name == "" {
			return errors.New("Required Name")
		}
		if t.StaticName == "" {
			return errors.New("Required StaticName")
		}
		if t.Subject == "" {
			return errors.New("Required Subject")
		}
		if t.Body == "" {
			return errors.New("Required Body")
		}
		if int64(strings.Count(t.Body, "$")) != t.ParamsCount {
			return errors.New("Parameters required don't match")
		}
		return nil

	case "add":
		if t.Name == "" {
			return errors.New("Required Name")
		}
		if t.StaticName == "" {
			return errors.New("Required StaticName")
		}
		if t.Subject == "" {
			return errors.New("Required Subject")
		}
		if t.Body == "" {
			return errors.New("Required Body")
		}
		if int64(strings.Count(t.Body, "$")) != t.ParamsCount {
			return errors.New("Parameters required don't match")
		}
		return nil
	default:
		return errors.New("Invalid validation operation")
	}
}

// PrepareTemplateOutput prepare the template to output
func (t *Template) PrepareTemplateOutput(temp Template) {
	t.ID = temp.ID
	t.Name = temp.Name
	t.StaticName = temp.StaticName
	t.Subject = temp.Subject
	t.Body = temp.Body
	t.ParamsCount = temp.ParamsCount
}
