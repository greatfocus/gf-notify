package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
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
func (repo *TemplateRepository) AddTemplate(template models.Template) (models.Template, error) {
	statement := `
    insert into template (name, staticName, subject, body, paramsCount)
    values ($1, $2, $3, $4, $5)
    returning id
  `
	var id int64
	err := repo.db.Select(statement, template.Name, template.StaticName, template.Subject, template.Body, template.ParamsCount).Scan(&id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return template, err
	}
	template.ID = id
	repo.deleteCache()
	return template, nil
}

// UpdateTemplate makes changes to the template
func (repo *TemplateRepository) UpdateTemplate(template models.Template) error {
	query := `
    update template
	set 
		name=$2,
		subject=$3,
		body=$4,
		paramsCount=$5,
		updatedOn=CURRENT_TIMESTAMP,
		enabled=$6
    where id=$1 and deleted=false
  	`
	res, err := repo.db.Update(query, template.ID, template.Name, template.Subject, template.Body, template.ParamsCount, template.Enabled)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated template for %d", template.ID)
	}

	repo.deleteCache()
	return nil
}

// GetTemplates method returns templates from the database
func (repo *TemplateRepository) GetTemplates(page int64) ([]models.Template, error) {
	// get data from cache
	var key = "TemplateRepository.GetTemplates" + strconv.Itoa(int(page))
	found, cache := repo.getTemplatesCache(key)
	if found {
		return cache, nil
	}

	query := `
	select id, name, staticName, subject, body, paramsCount, createdOn, updatedOn, enabled 
	from template 
	where deleted = false
	order by id DESC limit 500 OFFSET $1-1
	`
	rows, err := repo.db.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := templateMapper(rows)
	if err != nil {
		return nil, err
	}

	// update cache
	repo.setTemplatesCache(key, result)
	return result, nil
}

// GetTemplate method returns template from the database
func (repo *TemplateRepository) GetTemplate(id int64) (models.Template, error) {
	// get data from cache
	var key = "TemplateRepository.GetTemplate" + strconv.Itoa(int(id))
	found, cache := repo.getTemplateCache(key)
	if found {
		return cache, nil
	}

	template := models.Template{}
	query := `
	select id, name, staticName, subject, body, paramsCount, createdOn, updatedOn, enabled 
	from template 
	where id=$1 and deleted = false
	`
	row := repo.db.Select(query, id)
	err := row.Scan(&template.ID, &template.Name, &template.StaticName, &template.Subject, &template.Body, &template.ParamsCount, &template.CreatedOn, &template.UpdatedOn, &template.Enabled)
	if err != nil {
		return template, err
	}

	// update cache
	repo.setTemplateCache(key, template)
	return template, nil
}

// DeleteTemplate makes changes to the template delete status
func (repo *TemplateRepository) DeleteTemplate(id int64) error {
	query := `
    update template
	set 
		staticName=CONCAT(id, '-', staticName, 'DELETED'),
		updatedOn=CURRENT_TIMESTAMP,
		enabled=false,
		deleted=true
    where id=$1
  	`
	res, err := repo.db.Update(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated template for %d", id)
	}

	repo.deleteCache()
	return nil
}

// prepare template row
func templateMapper(rows *sql.Rows) ([]models.Template, error) {
	templates := []models.Template{}
	for rows.Next() {
		var template models.Template
		err := rows.Scan(&template.ID, &template.Name, &template.StaticName, &template.Subject, &template.Body, &template.ParamsCount, &template.CreatedOn, &template.UpdatedOn, &template.Enabled)
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
