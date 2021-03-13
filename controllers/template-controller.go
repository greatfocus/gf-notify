package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/response"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// TemplateController struct
type TemplateController struct {
	templateRepository *repositories.TemplateRepository
}

// Init method
func (t *TemplateController) Init(db *database.Conn, cache *cache.Cache) {
	t.templateRepository = &repositories.TemplateRepository{}
	t.templateRepository.Init(db, cache)
}

// Handler method routes to http methods supported
func (t *TemplateController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t.get(w, r)
	case http.MethodPost:
		t.add(w, r)
	case http.MethodPut:
		t.update(w, r)
	case http.MethodDelete:
		t.delete(w, r)
	default:
		err := errors.New("Invalid Request")
		response.Error(w, http.StatusNotFound, err)
		return
	}
}

// requestMessage method get templates
func (t *TemplateController) get(w http.ResponseWriter, r *http.Request) {
	pageStr := r.FormValue("page")

	if len(pageStr) != 0 {
		page, err := strconv.ParseInt(pageStr, 10, 64)
		templates := []models.Template{}
		templates, err = t.templateRepository.GetTemplates(page)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		response.Success(w, http.StatusOK, templates)
		return
	}

	derr := errors.New("Invalid parameter")
	response.Error(w, http.StatusBadRequest, derr)
	return
}

// add method adds new template
func (t *TemplateController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	template := models.Template{}
	err = json.Unmarshal(body, &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	template.PrepareTempate()
	err = template.ValidateTemplate("add")
	if err != nil {
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdTemplate, err := t.templateRepository.AddTemplate(template)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(createdTemplate)
	response.Success(w, http.StatusOK, result)
	return
}

// update method adds new template
func (t *TemplateController) update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}
	template := models.Template{}
	err = json.Unmarshal(body, &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusBadGateway, derr)
		return
	}

	err = template.ValidateTemplate("edit")
	if err != nil {
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = t.templateRepository.UpdateTemplate(template)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		response.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(template)
	response.Success(w, http.StatusOK, result)
	return
}

// requestMessage method delete templates
func (t *TemplateController) delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		err = t.templateRepository.DeleteTemplate(id)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		response.Success(w, http.StatusOK, id)
		return
	}

	derr := errors.New("Invalid parameter")
	response.Error(w, http.StatusBadRequest, derr)
	return
}
