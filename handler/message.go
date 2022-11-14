package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-sframe/server"
)

// Message struct
type Message struct {
	ChannelHandler    func(http.ResponseWriter, *http.Request)
	messageRepository *repositories.MessageRepository
	meta              *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hdl.create(w, r)
		return
	default:
		// catch all
		// if no method is satisfied return an error
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Allow", "GET")
	}
}

// ValidateRequest checks if request is valid
func (hdl Message) ValidateRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := hdl.meta.Request(w, r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Init method
func (hdl *Message) Init(meta *server.Meta) {
	hdl.messageRepository = &repositories.MessageRepository{}
	hdl.messageRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// add method adds new message
func (hdl *Message) create(w http.ResponseWriter, r *http.Request) {
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
	err = message.Validate("new")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	createdMessage, err := hdl.messageRepository.Create(ctx, hdl.meta.JWT.Secret, "staging", message)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	result := message.PrepareOutput(createdMessage)
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, result)
}
