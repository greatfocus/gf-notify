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

// Template struct
type Template struct {
	ChannelHandler     func(http.ResponseWriter, *http.Request)
	templateRepository *repositories.TemplateRepository
	meta               *server.Meta
}

// ServeHTTP checks if is valid method
func (hdl Template) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (hdl Template) ValidateRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := hdl.meta.Request(w, r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Init method
func (hdl *Template) Init(meta *server.Meta) {
	hdl.templateRepository = &repositories.TemplateRepository{}
	hdl.templateRepository.Init(meta.DB, meta.Cache)
	hdl.meta = meta
}

// add method adds new template
func (hdl *Template) create(w http.ResponseWriter, r *http.Request) {
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
	template := models.Template{}
	err = json.Unmarshal(buf.Bytes(), &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}
	template.PrepareTempate()
	err = template.ValidateTemplate("add")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	createdTemplate, err := hdl.templateRepository.Create(ctx, hdl.meta.JWT.Secret, template)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(createdTemplate)
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, result)
}

// requestMessage method get templates
func (hdl *Template) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	id := r.FormValue("id")
	if id != "" {
		template, err := hdl.templateRepository.GetTemplateByID(ctx, hdl.meta.JWT.Secret, id)
		if err != nil {
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, template)
	} else {
		lastID := r.FormValue("lastId")
		templates, err := hdl.templateRepository.GetTemplates(ctx, hdl.meta.JWT.Secret, lastID)
		if err != nil {
			hdl.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		hdl.meta.Success(w, r, templates)
	}
}

// update method adds new template
func (hdl *Template) update(w http.ResponseWriter, r *http.Request) {
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

	template := models.Template{}
	err = json.Unmarshal(buf.Bytes(), &template)
	if err != nil {
		derr := errors.New("invalid payload request")
		hdl.meta.Error(w, r, derr)
		return
	}

	err = template.ValidateTemplate("edit")
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	err = hdl.templateRepository.Update(ctx, hdl.meta.JWT.Secret, template)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}

	result := models.Template{}
	result.PrepareTemplateOutput(template)
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, result)
}

// requestMessage method delete templates
func (hdl *Template) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(hdl.meta.Timeout)*time.Second)
	defer cancel()

	id := r.FormValue("id")
	err := hdl.templateRepository.Delete(ctx, id)
	if err != nil {
		hdl.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	hdl.meta.Success(w, r, nil)
}
