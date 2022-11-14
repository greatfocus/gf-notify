package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/greatfocus/gf-notify/models"
	cache "github.com/greatfocus/gf-sframe/cache"
	"github.com/greatfocus/gf-sframe/database"
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

// Create a channel
func (repo *ChannelRepository) Create(ctx context.Context, enKey string, channel models.Channel) (models.Channel, error) {
	var id = uuid.New().String()
	statement := `
    insert into channel (id, name, key, url, user, pass)
    values ($1, PGP_SYM_ENCRYPT($2, '` + enKey + `'), PGP_SYM_ENCRYPT($3, '` + enKey + `'), PGP_SYM_ENCRYPT($4, '` + enKey + `'), PGP_SYM_ENCRYPT($5, '` + enKey + `'), PGP_SYM_ENCRYPT($6, '` + enKey + `'))
    returning id
  `
	_, inserted := repo.db.Insert(ctx, statement, id, channel.Name, channel.Key, channel.URL, channel.User, channel.Pass)
	if !inserted {
		return channel, errors.New("create channel failed")
	}
	channel.ID = id
	repo.deleteCache()
	return channel, nil
}

// GetChannels method returns channels from the database
func (repo *ChannelRepository) GetChannels(ctx context.Context, enKey string, lastID string) ([]models.Channel, error) {
	// get data from cache
	var key = "ChannelRepository.GetChannels"
	found, cache := repo.getChannelsCache(key)
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
			pgp_sym_decrypt(user::bytea, '` + enKey + `'),
			createdOn,
			updatedOn,
			enabled 
		from channel
		where id >= $1
		order BY createdOn DESC limit 20
		`
		rows, err = repo.db.Query(ctx, statement, lastID)
	} else {
		statement = `
		select
			id,
			pgp_sym_decrypt(name::bytea, '` + enKey + `'),
			pgp_sym_decrypt(key::bytea, '` + enKey + `'),
			pgp_sym_decrypt(user::bytea, '` + enKey + `'),
			createdOn,
			updatedOn,
			enabled 
		from channel 
		order BY createdOn DESC
		`
		rows, err = repo.db.Query(ctx, statement)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	result, err := channelMapper(rows)
	if err != nil {
		return nil, err
	}

	// update cache
	repo.setChannelsCache(key, result)
	return result, nil
}

// GetChannels method returns channels from the database
func (repo *ChannelRepository) GetChannelByID(ctx context.Context, enKey string, id string) (models.Channel, error) {
	// get data from cache
	var key = "ChannelRepository.GetChannels" + id
	found, cache := repo.getChannelCache(key)
	if found {
		return cache, nil
	}

	statement := `
	select
		id,
		pgp_sym_decrypt(name::bytea, '` + enKey + `'),
		pgp_sym_decrypt(key::bytea, '` + enKey + `'),
		pgp_sym_decrypt(user::bytea, '` + enKey + `'),
		createdOn,
		updatedOn,
		enabled 
	from channel
	where id = $1
	`
	row := repo.db.Select(ctx, statement, id)
	channel := models.Channel{}
	err := row.Scan(&channel.ID, &channel.Name, &channel.Key, channel.User, channel.CreatedOn, &channel.UpdatedOn, channel.Enabled)
	switch err {
	case sql.ErrNoRows:
		return channel, err
	case nil:
		// update cache
		repo.setChannelCache(key, channel)
		return channel, nil
	default:
		return channel, err
	}
}

// UpdateChannel makes changes to the channel
func (repo *ChannelRepository) Update(ctx context.Context, enKey string, channel models.Channel) error {
	statement := `
    UPDATE channel
	SET 
		name=PGP_SYM_ENCRYPT($2, '` + enKey + `'),
		url=PGP_SYM_ENCRYPT($3, '` + enKey + `'),
		user=PGP_SYM_ENCRYPT($4, '` + enKey + `'),
		pass=PGP_SYM_ENCRYPT($5, '` + enKey + `'),
		enabled=PGP_SYM_ENCRYPT($6, '` + enKey + `')
    WHERE id=$1
  	`
	updated := repo.db.Update(ctx, statement, channel.ID, channel.Name, channel.URL, channel.User, channel.Pass, channel.Enabled)
	if !updated {
		return errors.New("update channel failed")
	}

	repo.deleteCache()
	return nil
}

// Delete method
func (repo *ChannelRepository) Delete(ctx context.Context, id string) error {
	query := `
    delete from channel
    where id=$1
  	`
	deleted := repo.db.Delete(ctx, query, id)
	if !deleted {
		return errors.New("update channel failed")
	}

	repo.deleteCache()
	return nil
}

// GetChannelByKey method returns channels from the database
func (repo *ChannelRepository) GetChannelByKey(ctx context.Context, enKey string, key string) (models.Channel, error) {
	// get data from cache
	var cacheKey = "ChannelRepository.GetChannels" + key
	found, cache := repo.getChannelCache(cacheKey)
	if found {
		return cache, nil
	}

	statement := `
	select
		id,
		pgp_sym_decrypt(name::bytea, '` + enKey + `'),
		pgp_sym_decrypt(key::bytea, '` + enKey + `'),
		pgp_sym_decrypt(user::bytea, '` + enKey + `'),
		createdOn,
		updatedOn,
		enabled 
	from channel
	where key = $1
	`
	row := repo.db.Select(ctx, statement, key)
	channel := models.Channel{}
	err := row.Scan(&channel.ID, &channel.Name, &channel.Key, channel.User, channel.CreatedOn, &channel.UpdatedOn, channel.Enabled)
	switch err {
	case sql.ErrNoRows:
		return channel, err
	case nil:
		// update cache
		repo.setChannelCache(cacheKey, channel)
		return channel, nil
	default:
		return channel, err
	}
}

// prepare row
func channelMapper(rows *sql.Rows) ([]models.Channel, error) {
	channels := []models.Channel{}
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Key, &channel.User, &channel.CreatedOn, &channel.UpdatedOn, &channel.Enabled)
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

// getChannelCache method get cache for channel
func (repo *ChannelRepository) getChannelCache(key string) (bool, models.Channel) {
	var data models.Channel
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Channel)
		return found, data
	}
	return false, data
}

// setChannelCache method set cache for channel
func (repo *ChannelRepository) setChannelCache(key string, channel models.Channel) {
	if channel != (models.Channel{}) {
		channelRepositoryCacheKeys = append(channelRepositoryCacheKeys, key)
		repo.cache.Set(key, channel, 5*time.Minute)
	}
}
