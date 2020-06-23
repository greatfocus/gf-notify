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

// MessageRepository struct
type MessageRepository struct {
	db *database.DB
}

// Init method
func (repo *MessageRepository) Init(db *database.DB) {
	repo.db = db
}

// RequestMessage method created new message request
func (repo *MessageRepository) RequestMessage(message models.Message) (models.Message, error) {
	year, month, _ := time.Now().Date()
	statement := `
    insert into $table (channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    returning id
	`
	query := getDatabase(statement, year, int(month))
	var id int64
	err := repo.db.Conn.QueryRow(query, message.ChannelID, message.Recipient, message.Subject, message.Content, message.CreatedBy,
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
func (repo *MessageRepository) GetMessages(channelID int64, page int64, year int, month int) ([]models.Message, error) {

	statement := `
	select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdBy, m.createdOn, m.expireOn, m.statusId, s.name as status, m.attempts, m.priority 
	from $table m
	inner join channel c on c.id = m.channelId
	inner join status s on s.id = m.statusId
	where channelId = $1
	order BY createdOn ASC limit 50 OFFSET $2-1
	`
	query := getDatabase(statement, year, month)
	rows, err := repo.db.Conn.Query(query, channelID, page)
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
		err := rows.Scan(&message.ID, &message.ChannelID, &message.Channel, &message.Recipient, &message.Subject, &message.Content,
			&message.CreatedBy, &message.CreatedOn, &message.ExpireOn, &message.StatusID, &message.Status,
			&message.Attempts, &message.Priority)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
