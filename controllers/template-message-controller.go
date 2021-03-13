package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/response"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// TemplateMessageController struct
type TemplateMessageController struct {
	messageRepository  *repositories.MessageRepository
	templateRepository *repositories.TemplateRepository
}

// Init method
func (t *TemplateMessageController) Init(db *database.Conn, cache *cache.Cache) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(db, cache)
	t.templateRepository = &repositories.TemplateRepository{}
	t.templateRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (t *TemplateMessageController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.add(w, r)
	default:
		err := errors.New("Invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// add method adds new message
func (t *TemplateMessageController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	message := models.Message{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	message.PrepareInput(r)
	err = message.Validate("new-template")
	if err != nil {
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	newMessage, err := createTemplateMessage(t.messageRepository, t.templateRepository, message)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	createdMessage, err := t.messageRepository.Add("staging", newMessage)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Message{}
	result.PrepareOutput(createdMessage)
	response.Success(w, http.StatusOK, result)
	return
}

func createTemplateMessage(messageRepo *repositories.MessageRepository, templateRepo *repositories.TemplateRepository, message models.Message) (models.Message, error) {
	template, err := getTemplate(message.TemplateID, templateRepo)
	if err != nil {
		return message, err
	}

	if len(template.Name) < 1 {
		return message, errors.New("Template does not exist")
	}

	err = validateTemplate(message, template)
	if err != nil {
		return message, err
	}

	message.Content = createContent(template.Body, message.Params)
	message.Subject = template.Subject
	return message, nil
}

func getTemplate(id int64, repo *repositories.TemplateRepository) (models.Template, error) {
	template, err := repo.GetTemplate(id)
	if err != nil {
		return template, err
	}

	return template, nil
}

// validateTemplate checks if parameters are expected
func validateTemplate(message models.Message, template models.Template) error {
	if len(message.Params) != int(template.ParamsCount) {
		return errors.New("Parameters required don't match")
	}

	return nil
}

// createContent replaces params in the template
func createContent(body string, params []string) string {
	content := body
	for i := 0; i < len(params); i++ {
		temp := strings.Replace(content, "$"+strconv.Itoa(i+1), params[i], -1)
		content = temp
	}

	return content
}
