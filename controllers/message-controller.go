package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// MessageController struct
type MessageController struct {
	messageRepository *repositories.MessageRepository
}

// Init method
func (c *MessageController) Init(db *database.DB) {
	c.messageRepository = &repositories.MessageRepository{}
	c.messageRepository.Init(db)
}

// Handler method routes to http methods supported
func (c *MessageController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.add(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// add method adds new message
func (c *MessageController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	message := models.Message{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	err = message.PrepareInput(r)
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = message.Validate("new")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdMessage, err := c.messageRepository.Add("staging", message)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Message{}
	result.PrepareOutput(createdMessage)
	responses.Success(w, http.StatusCreated, result)
}
