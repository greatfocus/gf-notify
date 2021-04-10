package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/response"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// MessageController struct
type MessageController struct {
	messageRepository *repositories.MessageRepository
}

// Init method
func (c *MessageController) Init(db *database.Conn, cache *cache.Cache) {
	c.messageRepository = &repositories.MessageRepository{}
	c.messageRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (c *MessageController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.add(w, r)
	default:
		err := errors.New("invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// add method adds new message
func (c *MessageController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	message := models.Message{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	message.PrepareInput(r)
	err = message.Validate("new")
	if err != nil {
		log.Printf("error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdMessage, err := c.messageRepository.Add("staging", message)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Message{}
	result.PrepareOutput(createdMessage)
	response.Success(w, http.StatusOK, result)
}
