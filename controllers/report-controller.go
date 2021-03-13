package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/response"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// ReportController struct
type ReportController struct {
	messageRepository *repositories.MessageRepository
}

// Init method
func (d *ReportController) Init(db *database.Conn, cache *cache.Cache) {
	d.messageRepository = &repositories.MessageRepository{}
	d.messageRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (d *ReportController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.getReport(w, r)
	default:
		err := errors.New("Invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// getDashboard method gets dashboard
func (d *ReportController) getReport(w http.ResponseWriter, r *http.Request) {
	tableStr := r.FormValue("table")
	channelStr := r.FormValue("channel")
	yearStr := r.FormValue("year")
	monthStr := r.FormValue("month")
	pageStr := r.FormValue("page")
	if len(tableStr) != 0 && len(channelStr) != 0 && len(yearStr) != 0 && len(monthStr) != 0 && len(pageStr) != 0 {
		channel, err := strconv.ParseInt(channelStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			response.Error(w, http.StatusBadRequest, derr)
			return
		}

		year, err := strconv.ParseInt(yearStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			response.Error(w, http.StatusBadRequest, derr)
			return
		}

		month, err := strconv.ParseInt(monthStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			response.Error(w, http.StatusBadRequest, derr)
			return
		}

		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			response.Error(w, http.StatusBadRequest, derr)
			return
		}

		messages := []models.Message{}
		messages, err = d.messageRepository.ReportMessages(tableStr, channel, year, month, page)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		response.Success(w, http.StatusOK, messages)
		return
	}

	derr := errors.New("Invalid parameter")
	response.Error(w, http.StatusBadRequest, derr)
	return
}
