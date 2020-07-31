package repositories

import (
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// DashboardRepository struct
type DashboardRepository struct {
	db *database.DB
}

// Init method
func (repo *DashboardRepository) Init(db *database.DB) {
	repo.db = db
}

// GetDashboard method returns dashboard from the database
func (repo *DashboardRepository) GetDashboard(year int64, month int64) (models.Dashboard, error) {
	dashboard := models.Dashboard{}
	query := `
	select id, request, staging, queue, complete, failed
	from dashboard 
	WHERE year = $1 AND month = $2;
	`
	row := repo.db.Conn.QueryRow(query, year, month)
	err := row.Scan(&dashboard.ID, &dashboard.Staging, &dashboard.Queue, &dashboard.Complete, &dashboard.Failed)
	if err != nil {
		return models.Dashboard{}, err
	}

	return dashboard, nil
}
