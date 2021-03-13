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

// DashboardController struct
type DashboardController struct {
	dashboardRepository *repositories.DashboardRepository
}

// Init method
func (d *DashboardController) Init(db *database.Conn, cache *cache.Cache) {
	d.dashboardRepository = &repositories.DashboardRepository{}
	d.dashboardRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (d *DashboardController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.getDashboard(w, r)
	default:
		err := errors.New("Invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// getDashboard method gets dashboard
func (d *DashboardController) getDashboard(w http.ResponseWriter, r *http.Request) {
	yearStr := r.FormValue("year")
	monthStr := r.FormValue("month")
	if len(yearStr) != 0 && len(monthStr) != 0 {
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

		dashboard := models.Dashboard{}
		dashboard, err = d.dashboardRepository.GetDashboard(year, month)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		response.Success(w, http.StatusOK, dashboard)
		return
	}

	derr := errors.New("Invalid parameter")
	response.Error(w, http.StatusBadRequest, derr)
	return
}
