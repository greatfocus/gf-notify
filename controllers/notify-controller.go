package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// NotifyController struct
type NotifyController struct {
	userRepository *repositories.UserRepository
}

// Init method
func (c *NotifyController) Init(db *database.DB) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)
}

// Handler method routes to http methods supported
func (c *NotifyController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.getUsers(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// getUsers method
func (c *NotifyController) getUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		_, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		user := models.User{}
		//user, err := c.userRepository.GetUser(id)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, user)
		return
	}
	if len(pageStr) != 0 {
		_, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		users := []models.User{}
		//users, err = c.userRepository.GetUsers(page)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, users)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}
