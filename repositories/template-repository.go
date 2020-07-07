package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// TemplateRepository struct
type TemplateRepository struct {
	db *database.DB
}

// Init method
func (repo *TemplateRepository) Init(db *database.DB) {
	repo.db = db
}

// AddTemplate makes changes to the template
func (repo *TemplateRepository) AddTemplate(template models.Template) (models.Template, error) {
	statement := `
    insert into template (name, staticName, subject, body, paramsCount, createdBy, updatedBy)
    values ($1, $2, $3, $4, $5, $6, $7)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, template.Name, template.StaticName, template.Subject, template.Body, template.ParamsCount, template.CreatedBy, template.UpdatedBy).Scan(&id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return template, err
	}
	template.ID = id
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
		updatedBy=$6,
		enabled=$7
    where id=$1
  	`
	res, err := repo.db.Conn.Exec(query, template.ID, template.Name, template.Subject, template.Body, template.ParamsCount, template.UpdatedBy, template.Enabled)
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

	return nil
}

// GetTemplates method returns templates from the database
func (repo *TemplateRepository) GetTemplates(page int64) ([]models.Template, error) {
	query := `
	select id, name, staticName, subject, body, paramsCount, createdOn, updatedOn, enabled 
	from template 
	order by createdOn asc limit 500 OFFSET $1-1
	`
	rows, err := repo.db.Conn.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return templateMapper(rows)
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
