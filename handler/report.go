package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-sframe/server"
)

// Report struct
type Report struct {
	FileHandler       func(http.ResponseWriter, *http.Request)
	messageRepository *repositories.MessageRepository
	meta              *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl Report) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		hdl.get(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET")
}

// Init method
func (hdl *Report) Init(meta *server.Meta) {
	hdl.messageRepository = &repositories.MessageRepository{}
	hdl.messageRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// getDashboard method gets dashboard
func (hdl *Report) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	table := r.FormValue("table")
	channelID := r.FormValue("channelId")
	year := r.FormValue("year")
	month := r.FormValue("month")
	lastID := r.FormValue("page")
	if table != "" && channelID != "" && year != "" && month != "" && lastID != "" {
		messages, err := hdl.messageRepository.Report(ctx, hdl.meta.JWT.Secret, table, channelID, year, month, lastID)
		if err != nil {
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, messages)
		return
	}

	derr := errors.New("invalid payload request")
	w.WriteHeader(http.StatusUnprocessableEntity)
	hdl.meta.Error(w, r, derr)
}
