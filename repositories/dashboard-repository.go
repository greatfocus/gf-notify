package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/greatfocus/gf-notify/models"
	cache "github.com/greatfocus/gf-sframe/cache"
	"github.com/greatfocus/gf-sframe/database"
)

// dashboardRepositoryCacheKeys array
var dashboardRepositoryCacheKeys = []string{}

// DashboardRepository struct
type DashboardRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *DashboardRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// GetDashboard method returns dashboard from the database
func (repo *DashboardRepository) Get(ctx context.Context, year string, month string) (models.Dashboard, error) {
	// get data from cache
	var key = "DashboardRepository.GetDashboard" + year + month
	found, cache := repo.getDashboardCache(key)
	if found {
		return cache, nil
	}

	dashboard := models.Dashboard{}
	statement := `
	select id, request, staging, queue, complete, failed
	from dashboard 
	WHERE year = $1 AND month = $2;
	`
	row := repo.db.Select(ctx, statement, year, month)
	err := row.Scan(&dashboard.ID, &dashboard.Request, &dashboard.Staging, &dashboard.Queue, &dashboard.Complete, &dashboard.Failed)
	switch err {
	case sql.ErrNoRows:
		return dashboard, err
	case nil:
		// update cache
		repo.setDashboardCache(key, dashboard)
		return dashboard, nil
	default:
		return dashboard, err
	}
}

// updateStagingDashboard method make changes to dashboard
func (repo *DashboardRepository) Update(ctx context.Context, count int64) error {
	statement := `
	UPDATE dashboard
	SET 
		request = request + $1,
		staging = staging + $1
	WHERE year = (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP))
		AND month = (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP))
	`

	updated := repo.db.Update(ctx, statement, count)
	if !updated {
		return errors.New("update failed")
	}
	return nil
}

// getDashboardCache method get cache for dashboard
func (repo *DashboardRepository) getDashboardCache(key string) (bool, models.Dashboard) {
	var data models.Dashboard
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Dashboard)
		return found, data
	}
	return false, data
}

// setDashboardCache method set cache for dashboard
func (repo *DashboardRepository) setDashboardCache(key string, dashboard models.Dashboard) {
	if dashboard != (models.Dashboard{}) {
		dashboardRepositoryCacheKeys = append(dashboardRepositoryCacheKeys, key)
		repo.cache.Set(key, dashboard, 10*time.Minute)
	}
}
