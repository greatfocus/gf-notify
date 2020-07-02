package repositories

import (
	"context"
	"database/sql"
	"errors"
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
	Table     string
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
	Attempts  int64
	Page      int64
}

// Init method
func (repo *MessageRepository) Init(db *database.DB) {
	repo.db = db
}

/**
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
**/
func getTable(table string, query string) string {
	year, month, _ := time.Now().Date()
	yr := strconv.Itoa(year)
	mnth := strconv.Itoa(int(month))
	tableName := table + yr + mnth
	return strings.Replace(query, "$table", tableName, -1)
}

/**
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
**/
func getTableTime(query string) string {
	year, month, _ := time.Now().Date()
	yr := strconv.Itoa(year)
	mnth := strconv.Itoa(int(month))
	tt := yr + mnth
	return strings.Replace(query, "$tt", tt, -1)
}

// insert adds new  database record
func insert(repo *MessageRepository, query string, args []interface{}) (int64, error) {
	var id int64
	err := repo.db.Conn.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// update makes changes to database record
func update(repo *MessageRepository, query string, args []interface{}) (bool, error) {
	res, err := repo.db.Conn.Exec(query, args...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if count != 1 {
		err = errors.New("more than 1 record got updated")
		return false, err
	}

	return true, nil
}

// messageMapper prepare message row to object
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
	query := getTable(params.Table, params.Statement)
	rows, err := params.Repo.db.Conn.Query(query, params.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return messageMapper(rows)
}

// newQueryParams returns params
func newQueryParams(table string, statement string, args []interface{}, repo *MessageRepository) *QueryParam {
	// get current month and date
	year, month := getYearAndMonth()
	return &QueryParam{
		Table:     table,
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
func (repo *MessageRepository) GetMessages(table string, messageParam MessageParam) ([]models.Message, error) {
	// prepare the statement
	statement := `
	select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdBy, m.createdOn, m.expireOn, m.statusId, s.name as status, m.attempts, m.priority 
	from $table m
	inner join channel c on c.id = m.channelId
	inner join status s on s.id = m.statusId
	where 
		channelId = $1 and m.statusId=$2 and m.attempts<$3
	order BY createdOn ASC limit 500 OFFSET $4-1
	`
	// prepare the args
	var args []interface{}
	args = append(args, messageParam.ChannelID, messageParam.StatusID, messageParam.Attempts, messageParam.Page)
	queryParams := newQueryParams(table, statement, args, repo)
	return queryMessages(queryParams)
}

// Add method created new message request
func (repo *MessageRepository) Add(table string, message models.Message) (models.Message, error) {
	statement := `
    insert into $table (channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, priority)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
	`

	var args []interface{}
	args = append(args, message.ChannelID, message.Recipient, message.Subject, message.Content,
		message.CreatedBy, message.CreatedOn, message.ExpireOn, message.StatusID, message.Priority)
	query := getTable(table, statement)
	id, err := insert(repo, query, args)
	if err != nil {
		return message, err
	}

	message.ID = id
	return message, nil
}

// Update method make changes to message
func (repo *MessageRepository) Update(table string, message models.Message) (bool, error) {
	statement := `
    update $table
	set 
		statusId=$2, 
		attempts=$3,
		reference=$4	
    where id=$1
	`

	var args []interface{}
	args = append(args, message.ID, message.StatusID, message.Attempts, message.Reference)
	query := getTable(table, statement)
	success, err := update(repo, query, args)
	if err != nil {
		return success, err
	}

	return success, nil
}

// PopNewQueueToProcess runs a script to move data
func (repo *MessageRepository) PopNewQueueToProcess(table string, messageParam MessageParam) ([]models.Message, error) {
	statement := `
	WITH cte AS (
		select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdBy, m.createdOn, m.expireOn, m.statusId, s.name as status, m.attempts, m.priority 
		from $table m
		inner join channel c on c.id = m.channelId
		inner join status s on s.id = m.statusId
		where 
			m.channelId = $1 and m.statusId=$2 and m.attempts<$3
		order BY m.createdOn ASC limit 500 OFFSET $4-1
	)	
	update $table q
	set statusId = 3
	from cte
	where q.id=cte.id];
	`
	// prepare the args
	var args []interface{}
	args = append(args, messageParam.ChannelID, messageParam.StatusID, messageParam.Attempts, messageParam.Page)
	queryParams := newQueryParams(table, statement, args, repo)
	return queryMessages(queryParams)
}

// MoveStagedToQueue runs a script to move data
func (repo *MessageRepository) MoveStagedToQueue() (bool, error) {
	success := true
	statement := `
    WITH moved_rows AS (
		DELETE FROM staging$tt
		WHERE ID IN (
			SELECT ID
			FROM staging$tt
			ORDER BY createdOn
			LIMIT 500)
		RETURNING *
	)
	INSERT INTO queue$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority)
	SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, 2, 0, priority FROM moved_rows;
	`

	query := getTableTime(statement)
	_, err := repo.db.Conn.Exec(query)
	if err != nil {
		return false, err
	}

	return success, nil
}

// MoveOutFailedQueue runs a script to move data
func (repo *MessageRepository) MoveOutFailedQueue() (bool, error) {
	success := true
	statement := `
    WITH moved_rows AS (
		DELETE FROM queue$tt
		WHERE ID IN (
			SELECT ID
			FROM queue$tt
			WHERE attempts >= 5
			ORDER BY createdOn
			LIMIT 500)
		RETURNING *
	)
	INSERT INTO failed$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference)
	SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference FROM moved_rows;
	`

	query := getTableTime(statement)
	_, err := repo.db.Conn.Exec(query)
	if err != nil {
		return false, err
	}

	return success, nil
}

// MoveOutCompleteQueue runs a script to move data
func (repo *MessageRepository) MoveOutCompleteQueue() (bool, error) {
	success := true
	statement := `
    WITH moved_rows AS (
		DELETE FROM queue$tt
		WHERE ID IN (
			SELECT ID
			FROM queue$tt
			WHERE 
				attempts < 5
				AND statusId=4
			ORDER BY createdOn
			LIMIT 500)
		RETURNING *
	)
	INSERT INTO failed$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference)
	SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference FROM moved_rows;
	`

	query := getTableTime(statement)
	_, err := repo.db.Conn.Exec(query)
	if err != nil {
		return false, err
	}

	return success, nil
}
