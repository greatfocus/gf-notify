package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-sframe/cache"
	"github.com/greatfocus/gf-sframe/database"
)

// dashboardRepositoryCacheKeys array
var templateRepositoryCacheKeys = []string{}

// TemplateRepository struct
type TemplateRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *TemplateRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// AddTemplate makes changes to the template
func (repo *TemplateRepository) Create(ctx context.Context, enKey string, template models.Template) (models.Template, error) {
	var id = uuid.New().String()
	statement := ` xfbcf
    insert into template (id, name, key, subject, body)
    values ($1, PGP_SYM_ENCRYPT($2, '` + enKey + `'), PGP_SYM_ENCRYPT($3, '` + enKey + `'), PGP_SYM_ENCRYPT($4, '` + enKey + `'), PGP_SYM_ENCRYPT($5, '` + enKey + `'))
    returning id
  `
	_, inserted := repo.db.Insert(ctx, statement, id, template.Name, template.Key, template.Subject, template.Body)
	if !inserted {
		return template, errors.New("create template failed")
	}
	template.ID = id
	repo.deleteCache()
	return template, nil
}

// UpdateTemplate makes changes to the template
func (repo *TemplateRepository) Update(ctx context.Context, enKey string, template models.Template) error {
	query := `
    update template
	set 
		name=PGP_SYM_ENCRYPT($2, '` + enKey + `'),
		subject=PGP_SYM_ENCRYPT($3, '` + enKey + `'),
		body=PGP_SYM_ENCRYPT($4, '` + enKey + `'),
		updatedOn=CURRENT_TIMESTAMP,
		enabled=PGP_SYM_ENCRYPT($5, '` + enKey + `'),
    where id=$1 and deleted=false
  	`
	updated := repo.db.Update(ctx, query, template.ID, template.Name, template.Subject, template.Body, template.Enabled)
	if !updated {
		return errors.New("update template failed")
	}

	repo.deleteCache()
	return nil
}

// GetTemplates method returns templates from the database
func (repo *TemplateRepository) GetTemplates(ctx context.Context, enKey string, lastID string) ([]models.Template, error) {
	// get data from cache
	var key = "TemplateRepository.GetTemplates" + lastID
	found, cache := repo.getTemplatesCache(key)
	if found {
		return cache, nil
	}

	var statement string
	var rows *sql.Rows
	var err error
	if lastID != "" {
		statement = `
		select
			id,
			pgp_sym_decrypt(name::bytea, '` + enKey + `'),
			pgp_sym_decrypt(key::bytea, '` + enKey + `'),
			pgp_sym_decrypt(subject::bytea, '` + enKey + `'),
			pgp_sym_decrypt(body::bytea, '` + enKey + `'),
			createdOn, updatedOn, enabled 
		from template 
		where id >= $1 and deleted = false
		order BY createdOn DESC limit 20
		`
		rows, err = repo.db.Query(ctx, statement, lastID)
	} else {
		statement = `
		select
			id,
			pgp_sym_decrypt(name::bytea, '` + enKey + `'),
			pgp_sym_decrypt(key::bytea, '` + enKey + `'),
			pgp_sym_decrypt(subject::bytea, '` + enKey + `'),
			pgp_sym_decrypt(body::bytea, '` + enKey + `'),
			createdOn, updatedOn, enabled 
		from template 
		where deleted = false
		order BY createdOn DESC limit 20
		`
		rows, err = repo.db.Query(ctx, statement)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	result, err := templateMapper(rows)
	if err != nil {
		return nil, err
	}

	// update cache
	repo.setTemplatesCache(key, result)
	return result, nil
}

// GetTemplateByID method returns template from the database
func (repo *TemplateRepository) GetTemplateByID(ctx context.Context, enKey string, id string) (models.Template, error) {
	// get data from cache
	var key = "TemplateRepository.GetTemplate" + id
	found, cache := repo.getTemplateCache(key)
	if found {
		return cache, nil
	}

	template := models.Template{}
	query := `
	select
		id,
		pgp_sym_decrypt(name::bytea, '` + enKey + `'),
		pgp_sym_decrypt(key::bytea, '` + enKey + `'),
		pgp_sym_decrypt(subject::bytea, '` + enKey + `'),
		pgp_sym_decrypt(body::bytea, '` + enKey + `'),
		createdOn, updatedOn, enabled 
	from template 
	where id=$1 and deleted = false
	`
	row := repo.db.Select(ctx, query, id)
	err := row.Scan(&template.ID, &template.Name, &template.Key, &template.Subject, &template.Body, &template.CreatedOn, &template.UpdatedOn, &template.Enabled)
	if err != nil {
		return template, err
	}

	// update cache
	repo.setTemplateCache(key, template)
	return template, nil
}

// GetTemplate method returns template from the database
func (repo *TemplateRepository) GetTemplateByKey(ctx context.Context, enKey string, key string) (models.Template, error) {
	// get data from cache
	var cacheKey = "TemplateRepository.GetTemplate" + key
	found, cache := repo.getTemplateCache(cacheKey)
	if found {
		return cache, nil
	}

	template := models.Template{}
	query := `
	select
		id,
		pgp_sym_decrypt(name::bytea, '` + enKey + `'),
		pgp_sym_decrypt(key::bytea, '` + enKey + `'),
		pgp_sym_decrypt(subject::bytea, '` + enKey + `'),
		pgp_sym_decrypt(body::bytea, '` + enKey + `'),
		createdOn, updatedOn, enabled 
	from template 
	where key=$1 and deleted = false
	`
	row := repo.db.Select(ctx, query, key)
	err := row.Scan(&template.ID, &template.Name, &template.Key, &template.Subject, &template.Body, &template.CreatedOn, &template.UpdatedOn, &template.Enabled)
	if err != nil {
		return template, err
	}

	// update cache
	repo.setTemplateCache(cacheKey, template)
	return template, nil
}

// DeleteTemplate makes changes to the template delete status
func (repo *TemplateRepository) Delete(ctx context.Context, id string) error {
	query := `
	delete from template
    where id=$1
  	`
	deleted := repo.db.Delete(ctx, query, id)
	if !deleted {
		return errors.New("update template failed")
	}

	repo.deleteCache()
	return nil
}

// prepare template row
func templateMapper(rows *sql.Rows) ([]models.Template, error) {
	templates := []models.Template{}
	for rows.Next() {
		var template models.Template
		err := rows.Scan(&template.ID, &template.Name, &template.Key, &template.Subject, &template.Body, &template.CreatedOn, &template.UpdatedOn, &template.Enabled)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// getTemplateCache method get cache for template
func (repo *TemplateRepository) getTemplateCache(key string) (bool, models.Template) {
	var data models.Template
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Template)
		return found, data
	}
	return false, data
}

// setTemplateCache method set cache for template
func (repo *TemplateRepository) setTemplateCache(key string, template models.Template) {
	if template != (models.Template{}) {
		templateRepositoryCacheKeys = append(templateRepositoryCacheKeys, key)
		repo.cache.Set(key, template, 10*time.Minute)
	}
}

// getTemplatesCache method get cache for template
func (repo *TemplateRepository) getTemplatesCache(key string) (bool, []models.Template) {
	var data []models.Template
	if x, found := repo.cache.Get(key); found {
		data = x.([]models.Template)
		return found, data
	}
	return false, data
}

// setTemplateCache method set cache for template
func (repo *TemplateRepository) setTemplatesCache(key string, templates []models.Template) {
	if len(templates) > 0 {
		templateRepositoryCacheKeys = append(templateRepositoryCacheKeys, key)
		repo.cache.Set(key, templates, 10*time.Minute)
	}
}

// deleteCache method to delete
func (repo *TemplateRepository) deleteCache() {
	if len(templateRepositoryCacheKeys) > 0 {
		for i := 0; i < len(templateRepositoryCacheKeys); i++ {
			repo.cache.Delete(templateRepositoryCacheKeys[i])
		}
		templateRepositoryCacheKeys = []string{}
	}
}
