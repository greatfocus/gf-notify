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
	case http.MethodGet:
		c.getMessages(w, r)
	case http.MethodPost:
		c.requestMessage(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// getMessages method returns messages
func (c *MessageController) requestMessage(w http.ResponseWriter, r *http.Request) {
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

	createdMessage, err := c.messageRepository.RequestMessage(message)
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

// requestMessage method creates a message request
func (c *MessageController) getMessages(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")
	yearStr := r.FormValue("year")
	monthStr := r.FormValue("month")
	channelStr := r.FormValue("channelId")

	if len(pageStr) != 0 && len(yearStr) != 0 && len(monthStr) != 0 && len(monthStr) != 0 && len(channelStr) != 0 {
		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		year, err := strconv.ParseInt(yearStr, 10, 36)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		month, err := strconv.ParseInt(monthStr, 10, 36)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		channelID, err := strconv.ParseInt(channelStr, 10, 36)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		messages := []models.Message{}
		messages, err = c.messageRepository.GetMessages(channelID, page, int(year), int(month))
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, messages)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}