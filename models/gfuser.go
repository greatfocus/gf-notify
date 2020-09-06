package models

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/jwt"
	"github.com/greatfocus/gf-frame/utils"
)

// GFUser struct
type GFUser struct {
	ID        int64     `json:"id,omitempty"`
	RelatedID int64     `json:"relatedId,omitempty"`
	Email     string    `json:"email,omitempty"`
	Key       string    `json:"key,omitempty"`
	CreatedBy int64     `json:"-"`
	CreatedOn time.Time `json:"-"`
	UpdatedBy int64     `json:"-"`
	UpdatedOn time.Time `json:"-"`
	Enabled   bool      `json:"enabled,omitempty"`
	Deleted   bool      `json:"-"`
}

// PrepareUser initiliazes the user request object
func (s *GFUser) PrepareUser() {
	s.UpdatedOn = time.Now()
	s.CreatedOn = time.Now()
	key := strconv.Itoa(int(utils.Srand(10)))
	s.Key = utils.HashAndSalt([]byte(key))

	// TODO:consider making API call to users
	s.CreatedBy = 1
	s.UpdatedBy = 1
}

// PrepareUserEdit initiliazes the user request object
func (s *GFUser) PrepareUserEdit() {
	s.UpdatedOn = time.Now()
	key := strconv.Itoa(int(utils.Srand(10)))
	s.Key = utils.HashAndSalt([]byte(key))

	// TODO:consider making API call to users
	s.UpdatedBy = 1
}

// ValidationPermission check permissions
func (s *GFUser) ValidationPermission(r *http.Request) error {
	userID, _, err := jwt.ExtractTokenID(r)
	if err != nil {
		return errors.New("Invalid token")
	}

	// TODO:consider making API call to users
	s.CreatedBy = userID
	s.UpdatedBy = userID
	return nil
}

// ValidateUser check if request is valid
func (s *GFUser) ValidateUser(action string) error {
	switch strings.ToLower(action) {
	case "edit":
		if s.ID == 0 {
			return errors.New("Required ID")
		}
		if s.Key == "" {
			return errors.New("Required Key")
		}
		return nil
	case "add":
		if s.RelatedID == 0 {
			return errors.New("Required RelatedID")
		}
		if s.Email == "" {
			return errors.New("Required Email")
		}
		if s.Key == "" {
			return errors.New("Required Key")
		}
		return nil
	default:
		return errors.New("Invalid validation operation")
	}
}

// PrepareUserOutput prepare the user to output
func (s *GFUser) PrepareUserOutput(user GFUser) {
	s.ID = user.ID
	s.RelatedID = user.RelatedID
	s.Email = user.Email
	s.Key = user.Key
}
