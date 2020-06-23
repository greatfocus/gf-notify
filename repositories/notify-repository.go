package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// NotifyRepository struct
type NotifyRepository struct {
	db *database.DB
}

// Init method
func (repo *NotifyRepository) Init(db *database.DB) {
	repo.db = db
}

// RequestMessage method created new message request
func (repo *NotifyRepository) RequestMessage(message models.Message) (models.Message, error) {
	year, month, _ := time.Now().Date()
	database := getDatabase(int64(year), int64(month))
	statement := `
    insert into $1 (channel, recipient, content, createdBy, createdOn, expireOn, statusId, attempts, priority)
    values ($2, $3, $4, $5, $6, $7, $8, $9, $10)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, database, message.Channel, message.Recipient, message.Content, message.CreatedBy,
		message.CreatedOn, message.ExpireOn, message.StatusID, message.Attempts, message.Priority).Scan(&id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return message, err
	}
	createdMessage := message
	createdMessage.ID = id
	return createdMessage, nil
}

// GetMessages method returns messages from the database
func (repo *NotifyRepository) GetMessages(channel string, page int64, year int64, month int64) ([]models.Message, error) {
	database := getDatabase(year, month)
	query := `
	select id, channel, recipient, content, createdBy, createdOn, expireOn, statusId, attempts, priority, refId 
	from $1 
	where channel = $2
	order BY createdOn ASC limit 50 OFFSET $3-1
	`
	rows, err := repo.db.Conn.Query(query, database, channel, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return messageMapper(rows)
}

// getDatabase returns database name depending with the year and month
func getDatabase(year int64, month int64) string {
	return "messageOut" + string(year) + string(month)
}

// prepare users row
func messageMapper(rows *sql.Rows) ([]models.Message, error) {
	messages := []models.Message{}
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.Channel)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
