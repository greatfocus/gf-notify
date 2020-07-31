package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// GFUserController struct
type GFUserController struct {
	gfuserRepository *repositories.GFUserRepository
}

// Init method
func (c *GFUserController) Init(db *database.DB) {
	c.gfuserRepository = &repositories.GFUserRepository{}
	c.gfuserRepository.Init(db)
}

// Handler method routes to http methods supported
func (c *GFUserController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.get(w, r)
	case http.MethodPost:
		c.add(w, r)
	case http.MethodPut:
		c.update(w, r)
	case http.MethodDelete:
		c.delete(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// requestMessage method get users
func (c *GFUserController) get(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")

	if len(pageStr) != 0 {
		page, err := strconv.ParseInt(pageStr, 10, 64)
		users := []models.GFUser{}
		users, err = c.gfuserRepository.GetUsers(page)
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

// add method adds new user
func (c *GFUserController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user := models.GFUser{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user.PrepareUser()
	err = user.ValidateUser("add")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdUser, err := c.gfuserRepository.AddUser(user)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.GFUser{}
	result.PrepareUserOutput(createdUser)
	responses.Success(w, http.StatusCreated, result)
}

// update method adds new user
func (c *GFUserController) update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user := models.GFUser{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	user.PrepareUserEdit()
	err = user.ValidateUser("edit")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = c.gfuserRepository.UpdateUser(user)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.GFUser{}
	result.PrepareUserOutput(user)
	responses.Success(w, http.StatusCreated, result)
}

// requestMessage method delete user
func (c *GFUserController) delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		err = c.gfuserRepository.DeleteUser(id, 1)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, id)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}
