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

// MessageBulkController struct
type MessageBulkController struct {
	messageRepository *repositories.MessageRepository
}

// Init method
func (m *MessageBulkController) Init(db *database.DB) {
	m.messageRepository = &repositories.MessageRepository{}
	m.messageRepository.Init(db)
}

// Handler method routes to http methods supported
func (m *MessageBulkController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.addMessage(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// addMessage method adds new message
func (m *MessageBulkController) addMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	messages := []models.Message{}
	err = json.Unmarshal(body, &messages)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// maximum bulk insert is 100
	if len(messages) > 100 {
		err := errors.New("Required Content")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	Validate(w, messages)
	PrepareInput(w, r, messages)
	BulkInsert(w, r, m.messageRepository, messages)
}

// Validate bulk messages
func Validate(w http.ResponseWriter, messages []models.Message) {
	for i := 0; i < len(messages); i++ {
		err := messages[i].Validate("new")
		if err != nil {
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
}

// PrepareInput bulk messages
func PrepareInput(w http.ResponseWriter, r *http.Request, messages []models.Message) {
	for i := 0; i < len(messages); i++ {
		err := messages[i].PrepareInput(r)
		if err != nil {
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
}

// BulkInsert bulk messages
func BulkInsert(w http.ResponseWriter, r *http.Request, repo *repositories.MessageRepository, messages []models.Message) {
	for i := 0; i < len(messages); i++ {
		result, err := repo.Add("staging", messages[i])
		if err != nil {
			messages[i].Operation = "failed"
		}
		messages[i].Operation = "success"
		messages[i].ID = result.ID
	}
	responses.Success(w, http.StatusCreated, messages)
}

/*
// addMessage method creates a message request
func (c *MessageBulkController) getMessages(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")
	statusStr := r.FormValue("statusId")
	channelStr := r.FormValue("channelId")

	if len(pageStr) != 0 && len(statusStr) != 0 && len(channelStr) != 0 {
		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		statusID, err := strconv.ParseInt(statusStr, 10, 36)
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
		params := repositories.MessageParam{
			ChannelID: channelID,
			StatusID:  statusID,
			Page:      page,
		}
		messages, err = c.messageRepository.GetMessages(params)
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
*/
