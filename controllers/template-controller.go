package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// TemplateController struct
type TemplateController struct {
	templateRepository *repositories.TemplateRepository
}

// Init method
func (t *TemplateController) Init(db *database.DB) {
	t.templateRepository = &repositories.TemplateRepository{}
	t.templateRepository.Init(db)
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
		responses.Error(w, http.StatusUnprocessableEntity, err)
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
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, templates)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}

// add method adds new template
func (t *TemplateController) add(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	template := models.Template{}
	err = json.Unmarshal(body, &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	template.PrepareTempate()
	err = template.ValidateTemplate("add")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdTemplate, err := t.templateRepository.AddTemplate(template)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(createdTemplate)
	responses.Success(w, http.StatusCreated, result)
}

// update method adds new template
func (t *TemplateController) update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	template := models.Template{}
	err = json.Unmarshal(body, &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	err = template.ValidateTemplate("edit")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = t.templateRepository.UpdateTemplate(template)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(template)
	responses.Success(w, http.StatusCreated, result)
}

// requestMessage method delete templates
func (t *TemplateController) delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err := strconv.ParseInt(idStr, 10, 64)
		err = t.templateRepository.DeleteTemplate(id, 1)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, id)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}
