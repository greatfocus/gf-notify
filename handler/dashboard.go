package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-sframe/server"
)

// Dashboard struct
type Dashboard struct {
	DashboardHandler    func(http.ResponseWriter, *http.Request)
	dashboardRepository *repositories.DashboardRepository
	meta                *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		hdl.getDashboard(w, r)
		return
	}
	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET")
}

// Init method
func (hdl *Dashboard) Init(meta *server.Meta) {
	hdl.dashboardRepository = &repositories.DashboardRepository{}
	hdl.dashboardRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// getDashboard method gets dashboard
func (hdl *Dashboard) getDashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	year := r.FormValue("year")
	month := r.FormValue("month")
	if year != "" && month != "" {
		dashboard, err := hdl.dashboardRepository.Get(ctx, year, month)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, dashboard)
		return
	}

	derr := errors.New("invalid payload request")
	w.WriteHeader(http.StatusUnprocessableEntity)
	hdl.meta.Error(w, r, derr)
}
