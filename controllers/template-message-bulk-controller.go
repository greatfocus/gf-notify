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

// TemplateMessageBulkController struct
type TemplateMessageBulkController struct {
	messageRepository  *repositories.MessageRepository
	templateRepository *repositories.TemplateRepository
}

// Init method
func (m *TemplateMessageBulkController) Init(db *database.DB) {
	m.messageRepository = &repositories.MessageRepository{}
	m.messageRepository.Init(db)
	m.templateRepository = &repositories.TemplateRepository{}
	m.templateRepository.Init(db)
}

// Handler method routes to http methods supported
func (m *TemplateMessageBulkController) Handler(w http.ResponseWriter, r *http.Request) {
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
func (m *TemplateMessageBulkController) addMessage(w http.ResponseWriter, r *http.Request) {
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

	msg, err := PrepareTemplateInput(m.messageRepository, m.templateRepository, messages)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	Validate(w, msg)
	PrepareInput(w, r, msg)
	BulkInsert(w, r, m.messageRepository, msg)
}

// PrepareTemplateInput bulk messages
func PrepareTemplateInput(messageRepo *repositories.MessageRepository, templateRepo *repositories.TemplateRepository, messages []models.Message) ([]models.Message, error) {
	newMessages := []models.Message{}
	for i := 0; i < len(messages); i++ {
		msg, err := createTemplateMessage(messageRepo, templateRepo, messages[i])
		if err != nil {
			return nil, err
		}
		newMessages = append(newMessages, msg)
	}

	return newMessages, nil
}