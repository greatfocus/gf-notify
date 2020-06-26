package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/models"
)

// MessageRepository struct
type MessageRepository struct {
	db  *database.DB
	ctx context.Context
}

// QueryParam struct
type QueryParam struct {
	Month     int
	Year      int
	Statement string
	Args      []interface{}
	Repo      *MessageRepository
}

// MessageParam struct
type MessageParam struct {
	ChannelID int64
	StatusID  int64
	Page      int64
}

// Init method
func (repo *MessageRepository) Init(db *database.DB) {
	repo.db = db
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

// queryMessages executes to get messages
func queryMessages(params *QueryParam) ([]models.Message, error) {
	query := getDatabase(params.Statement, params.Year, params.Month)
	rows, err := params.Repo.db.Conn.Query(query, params.Args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return messageMapper(rows)
}

// newQueryParams returns params
func newQueryParams(year int, month int, statement string, args []interface{}, repo *MessageRepository) *QueryParam {
	return &QueryParam{
		Month:     month,
		Year:      year,
		Statement: statement,
		Args:      args,
		Repo:      repo,
	}
}

// getYearAndMonth returns the year and month
func getYearAndMonth() (year int, month int) {
	date := time.Now()
	year = date.Year()
	month = int(date.Month())
	return year, month
}

// GetMessages method returns messages from the database
func (repo *MessageRepository) GetMessages(messageParam MessageParam) ([]models.Message, error) {
	// prepare the statement
	statement := `
	select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdBy, m.createdOn, m.expireOn, m.statusId, s.name as status, m.attempts, m.priority 
	from $table m
	inner join channel c on c.id = m.channelId
	inner join status s on s.id = m.statusId
	where channelId = $1 and m.statusId=$2
	order BY createdOn ASC limit 60 OFFSET $3-1
	`
	// prepare the args
	var args []interface{}
	args = append(args, messageParam.ChannelID, messageParam.StatusID, messageParam.Page)

	// get current month and date
	year, month := getYearAndMonth()
	queryParams := newQueryParams(year, month, statement, args, repo)
	return queryMessages(queryParams)
}

// AddMessage method created new message request
func (repo *MessageRepository) AddMessage(message models.Message) (models.Message, error) {
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
	message.ID = id
	return message, nil
}

// UpdateMessage make changes to message
func (repo *MessageRepository) UpdateMessage(message models.Message) error {
	tx, err := repo.db.Conn.BeginTx(repo.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	statement := `
    update $table
	set 
		statusId=$2, 
		attempts=$3,
		refId=$4	
    where id=$1
	`
	year, month := getYearAndMonth()
	query := getDatabase(statement, year, int(month))
	res, execErr := tx.ExecContext(repo.ctx, query, message.ID, message.StatusID, message.Attempts, message.RefID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("update failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		log.Printf("update failed: %v", execErr)
		return execErr
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if count != 1 {
		err := fmt.Errorf("more than 1 record got updated User for %d", message.ID)
		return err
	}

	return nil
}
