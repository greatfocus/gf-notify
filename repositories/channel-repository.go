package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// channelRepositoryCacheKeys array
var channelRepositoryCacheKeys = []string{}

// ChannelRepository struct
type ChannelRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *ChannelRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// UpdateChannel makes changes to the channel
func (repo *ChannelRepository) UpdateChannel(channel models.Channel) error {
	query := `
    UPDATE channel
	SET 
		priority=$2,
		enabled=$3
    WHERE id=$1
  	`
	res, err := repo.db.Update(query, channel.ID, channel.Priority, channel.Enabled)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated Channel for %d", channel.ID)
	}

	repo.deleteCache()
	return nil
}

// GetChannels method returns channels from the database
func (repo *ChannelRepository) GetChannels() ([]models.Channel, error) {
	// get data from cache
	var key = "ChannelRepository.GetChannels"
	found, cache := repo.getChannelsCache(key)
	if found {
		return cache, nil
	}

	query := `
	select id, name, staticName, priority, createdOn, updatedOn, enabled 
	from channel 
	order BY id DESC
	`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := channelMapper(rows)
	if err != nil {
		return nil, err
	}

	// update cache
	repo.setChannelsCache(key, result)
	return result, nil
}

// prepare row
func channelMapper(rows *sql.Rows) ([]models.Channel, error) {
	channels := []models.Channel{}
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(&channel.ID, &channel.Name, &channel.StaticName, &channel.Priority, &channel.CreatedOn, &channel.UpdatedOn, &channel.Enabled)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}

// getChannelsCache method get cache for channels
func (repo *ChannelRepository) getChannelsCache(key string) (bool, []models.Channel) {
	var data []models.Channel
	if x, found := repo.cache.Get(key); found {
		data = x.([]models.Channel)
		return found, data
	}
	return false, data
}

// setChannelsCache method set cache for channels
func (repo *ChannelRepository) setChannelsCache(key string, channels []models.Channel) {
	if len(channels) > 0 {
		channelRepositoryCacheKeys = append(channelRepositoryCacheKeys, key)
		repo.cache.Set(key, channels, 10*time.Minute)
	}
}

// deleteCache method to delete
func (repo *ChannelRepository) deleteCache() {
	if len(channelRepositoryCacheKeys) > 0 {
		for i := 0; i < len(channelRepositoryCacheKeys); i++ {
			repo.cache.Delete(channelRepositoryCacheKeys[i])
		}
		channelRepositoryCacheKeys = []string{}
	}
}
