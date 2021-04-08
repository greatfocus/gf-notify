package repositories

import (
	"strconv"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
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
func (repo *DashboardRepository) GetDashboard(year int64, month int64) (models.Dashboard, error) {
	// get data from cache
	var key = "DashboardRepository.GetDashboard" + strconv.Itoa(int(year)) + strconv.Itoa(int(month))
	found, cache := repo.getDashboardCache(key)
	if found {
		return cache, nil
	}

	dashboard := models.Dashboard{}
	query := `
	select id, request, staging, queue, complete, failed
	from dashboard 
	WHERE year = $1 AND month = $2;
	`
	row := repo.db.Select(query, year, month)
	err := row.Scan(&dashboard.ID, &dashboard.Request, &dashboard.Staging, &dashboard.Queue, &dashboard.Complete, &dashboard.Failed)
	if err != nil {
		return models.Dashboard{}, err
	}

	// update cache
	repo.setDashboardCache(key, dashboard)
	return dashboard, nil
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
