package repositories

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
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
	statement := `
    insert into $table (channel, recipient, content, createdBy, createdOn, expireOn, statusId, attempts, priority)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
	`
	query := getDatabase(statement, year, int(month))
	var id int64
	err := repo.db.Conn.QueryRow(query, message.Channel, message.Recipient, message.Content, message.CreatedBy,
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
func (repo *NotifyRepository) GetMessages(channel string, page int64, year int, month int) ([]models.Message, error) {

	statement := `
	select id, channel, recipient, content, createdBy, createdOn, expireOn, statusId, attempts, priority 
	from $table 
	where channel = $1
	order BY createdOn ASC limit 50 OFFSET $2-1
	`
	query := getDatabase(statement, year, month)
	rows, err := repo.db.Conn.Query(query, channel, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return messageMapper(rows)
}

// getDatabase returns database name depending with the year and month
func getDatabase(query string, year int, month int) string {
	yr := strconv.Itoa(year)
	mnth := strconv.Itoa(month)
	database := "messageOut" + yr + mnth
	return strings.Replace(query, "$table", database, -1)
}

// prepare users row
func messageMapper(rows *sql.Rows) ([]models.Message, error) {
	messages := []models.Message{}
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.ID, &message.Channel, &message.Recipient, &message.Content,
			&message.CreatedBy, &message.CreatedOn, &message.ExpireOn, &message.StatusID,
			&message.Attempts, &message.Priority)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
