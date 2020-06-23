package repositories

import (
	"database/sql"
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// ChannelRepository struct
type ChannelRepository struct {
	db *database.DB
}

// Init method
func (repo *ChannelRepository) Init(db *database.DB) {
	repo.db = db
}

// UpdateChannel makes changes to the channel
func (repo *ChannelRepository) UpdateChannel(channel models.Channel) error {
	query := `
    update channel
	set 
		updatedBy=$2,
		priority=$3
		enabled=$4
    where id=$1
  	`
	res, err := repo.db.Conn.Exec(query, channel.ID, channel.UpdatedBy, channel.Priority, channel.Enabled)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated User for %d", channel.ID)
	}

	return nil
}

// GetChannels method returns channels from the database
func (repo *ChannelRepository) GetChannels() ([]models.Channel, error) {
	query := `
	select id, name, staticName, priority, updateBy, updateOn, enabled 
	from channel 
	order BY id ASC
	`
	rows, err := repo.db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return channelMapper(rows)
}

// prepare users row
func channelMapper(rows *sql.Rows) ([]models.Channel, error) {
	channels := []models.Channel{}
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(&channel.ID, &channel.Name, &channel.StaticName, &channel.Priority, &channel.UpdatedBy, &channel.UpdatedOn)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	return channels, nil
}
