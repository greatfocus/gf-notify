package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-sframe/server"
)

// TemplateMessage struct
type TemplateMessage struct {
	ChannelHandler     func(http.ResponseWriter, *http.Request)
	messageRepository  *repositories.MessageRepository
	templateRepository *repositories.TemplateRepository
	meta               *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl TemplateMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hdl.create(w, r)
		return
	default:
		// catch all
		// if no method is satisfied return an error
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Allow", "POST")
	}
}

// ValidateRequest checks if request is valid
func (hdl TemplateMessage) ValidateRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := hdl.meta.Request(w, r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Init method
func (hdl *TemplateMessage) Init(meta *server.Meta) {
	hdl.templateRepository = &repositories.TemplateRepository{}
	hdl.templateRepository.Init(meta.DB, meta.Cache)
	hdl.messageRepository = &repositories.MessageRepository{}
	hdl.messageRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// add method adds new message
func (hdl *TemplateMessage) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	data, err := hdl.ValidateRequest(w, r)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(data)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}

	message := models.Message{}
	err = json.Unmarshal(buf.Bytes(), &message)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}
	message.PrepareInput(r)
	err = message.Validate("new-template")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	template, err := hdl.templateRepository.GetTemplateByID(ctx, hdl.meta.JWT.Secret, message.TemplateID)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}
	if len(template.Name) < 1 {
		hdl.meta.Error(w, r, errors.New("template does not exist"))
		return
	}

	message.Content = createContent(template.Body, message.Params)
	message.Subject = template.Subject
	createdMessage, err := hdl.messageRepository.Create(ctx, hdl.meta.JWT.Secret, "staging", message)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	result := models.Message{}
	result.PrepareOutput(createdMessage)
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, result)
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
