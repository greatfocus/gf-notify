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

// Channel struct
type Channel struct {
	ChannelHandler    func(http.ResponseWriter, *http.Request)
	channelRepository *repositories.ChannelRepository
	meta              *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl Channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hdl.create(w, r)
		return
	case http.MethodGet:
		hdl.get(w, r)
		return
	case http.MethodPut:
		hdl.update(w, r)
		return
	case http.MethodDelete:
		hdl.delete(w, r)
		return
	default:
		// catch all
		// if no method is satisfied return an error
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Allow", "POST, GET, PUT, DELETE")
	}
}

// ValidateRequest checks if request is valid
func (hdl Channel) ValidateRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := hdl.meta.Request(w, r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Init method
func (hdl *Channel) Init(meta *server.Meta) {
	hdl.channelRepository = &repositories.ChannelRepository{}
	hdl.channelRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// create method returns add channel
func (hdl *Channel) create(w http.ResponseWriter, r *http.Request) {
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
	channel := models.Channel{}
	err = json.Unmarshal(buf.Bytes(), &channel)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}

	channel.PrepareChannel(r)
	err = channel.ValidateChannel("create")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	channel, err = hdl.channelRepository.Create(ctx, hdl.meta.JWT.Secret, channel)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, channel)
}

// get method gets channel(s)
func (hdl *Channel) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	id := r.FormValue("id")
	if id != "" {
		channels, err := hdl.channelRepository.GetChannelByID(ctx, hdl.meta.JWT.Secret, id)
		if err != nil {
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, channels)
	} else {
		lastID := r.FormValue("lastId")
		channels, err := hdl.channelRepository.GetChannels(ctx, hdl.meta.JWT.Secret, lastID)
		if err != nil {
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, channels)
	}
}

// update method changes channel
func (hdl *Channel) update(w http.ResponseWriter, r *http.Request) {
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
	channel := models.Channel{}
	err = json.Unmarshal(buf.Bytes(), &channel)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}

	channel.PrepareChannel(r)
	err = channel.ValidateChannel("update")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	err = hdl.channelRepository.Update(ctx, hdl.meta.JWT.Secret, channel)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, channel)
}

// requestMessage method creates a message request
func (hdl *Channel) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	id := r.FormValue("id")
	err := hdl.channelRepository.Delete(ctx, id)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, nil)
}
