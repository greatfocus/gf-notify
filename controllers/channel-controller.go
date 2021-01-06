package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// ChannelController struct
type ChannelController struct {
	channelRepository *repositories.ChannelRepository
}

// Init method
func (c *ChannelController) Init(db *database.Conn, cache *cache.Cache) {
	c.channelRepository = &repositories.ChannelRepository{}
	c.channelRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (c *ChannelController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.getChannels(w, r)
	case http.MethodPost:
		c.updateChannel(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// getMessages method returns messages
func (c *ChannelController) updateChannel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	channel := models.Channel{}
	err = json.Unmarshal(body, &channel)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}

	channel.PrepareChannel(r)
	err = channel.ValidateChannel("update")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = c.channelRepository.UpdateChannel(channel)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	responses.Success(w, http.StatusOK, channel)
	return
}

// requestMessage method creates a message request
func (c *ChannelController) getChannels(w http.ResponseWriter, r *http.Request) {
	channels := []models.Channel{}
	channels, err := c.channelRepository.GetChannels()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, channels)
	return
}
