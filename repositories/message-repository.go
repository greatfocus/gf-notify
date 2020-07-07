package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	Month     string
	Year      string
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

// getYearAndMonth returns the year and month
func getYearAndMonth() (y string, m string) {
	date := time.Now()
	year := date.Year()
	month := int(date.Month())
	y = strconv.Itoa(year)
	m = strconv.Itoa(int(month))
	return y, m
}

/**
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
**/
func getTable(name string) string {
	year, month := getYearAndMonth()
	tableName := name + year + month
	return tableName
}

/**
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
**/
func replaceTimeHolder(query string) string {
	year, month := getYearAndMonth()
	tt := year + month
	return strings.Replace(query, "$tt", tt, -1)
}

// replaceTableHolder change holder with table name
func replaceTableHolder(name string, statement string) string {
	table := getTable(name)
	return strings.Replace(statement, "$table", table, -1)
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

// insert adds new  database record
func insert(params *QueryParam) (int64, error) {
	var id int64
	err := params.Repo.db.Conn.QueryRow(params.Statement, params.Args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// update makes changes to database record
func update(params *QueryParam, rowCount int64) (bool, error) {
	res, err := params.Repo.db.Conn.Exec(params.Statement, params.Args...)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if count != rowCount {
		err = errors.New("Records updated more than the expected")
		return false, err
	}

	return true, nil
}

// queryMessages executes to get messages
func query(params *QueryParam) ([]models.Message, error) {
	rows, err := params.Repo.db.Conn.Query(params.Statement, params.Args...)
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

// prepareInStatement for creating in statement
func prepareInStatement(arg []interface{}) string {
	statement := "("
	for i := 0; i < len(arg); i++ {
		statement = statement + fmt.Sprintf("%v", arg[i])
		record := (i + 1)
		if record == len(arg) {
			statement = statement + ")"
		} else {
			statement = statement + ","
		}
	}
	return statement
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
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo)
	return query(queryParams)
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
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo)
	id, err := insert(queryParams)
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
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo)
	success, err := update(queryParams, 1)
	if err != nil {
		return success, err
	}

	return success, nil
}

// UpdateQueueToProcessing method make changes to message
func (repo *MessageRepository) UpdateQueueToProcessing(table string, args []interface{}) (bool, error) {
	statement := `
    update queue$tt
	set 
		statusId=3
    where id IN $IN
	`
	newArgs := prepareInStatement(args)
	query := strings.Replace(statement, "$IN", newArgs, -1)
	qry := replaceTimeHolder(query)
	queryParams := newQueryParams(table, replaceTableHolder(table, qry), nil, repo)
	success, err := update(queryParams, int64(len(args)))
	if err != nil {
		return success, err
	}

	return success, nil
}

// MoveStagedToQueue runs a script to move data
func (repo *MessageRepository) MoveStagedToQueue() (bool, error) {
	success := true
	statement := `
	DO $$ 
	DECLARE
		mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
		yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));
	BEGIN
		WITH moved_rows AS (
			DELETE FROM staging$tt
			WHERE ID IN (
				SELECT ID
				FROM staging$tt
				ORDER BY createdOn
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO queue$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority)
			SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, 2, 0, priority FROM moved_rows
		)		
		UPDATE dashboard
		SET queue = queue + (SELECT COUNT(ID) FROM moved_rows), staging = staging + (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND mnth = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
	_, err := repo.db.Conn.Exec(query)
	if err != nil {
		return false, err
	}

	return success, nil
}

// ReQueueProcessingEmails runs a script to move data
func (repo *MessageRepository) ReQueueProcessingEmails() (bool, error) {
	success := true
	statement := `
		UPDATE queue$tt
		SET statusid=2
		WHERE statusid=3 and EXTRACT(MINUTE FROM updatedOn) > 10;
	`
	query := replaceTimeHolder(statement)
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
	DO $$ 
	DECLARE
		mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
		yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));

	BEGIN
		WITH moved_rows AS (
			DELETE FROM queue$tt
			WHERE ID IN (
				SELECT ID
				FROM queue$tt
				WHERE statusid=5
				ORDER BY createdOn
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO failed$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference)
			SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference FROM moved_rows
		)
		UPDATE dashboard
		SET 
			failed = failed + (SELECT COUNT(ID) FROM moved_rows),
			queue = queue - (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND mnth = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
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
	DO $$ 
	DECLARE
		mnth SMALLINT := (SELECT EXTRACT(MONTH FROM CURRENT_TIMESTAMP));
		yr SMALLINT := (SELECT EXTRACT(YEAR FROM CURRENT_TIMESTAMP));

	BEGIN
		WITH moved_rows AS (
			DELETE FROM queue$tt
			WHERE ID IN (
				SELECT ID
				FROM queue$tt
				WHERE statusId=4
				ORDER BY createdOn
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO complete$tt (id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference)
			SELECT id, channelId, recipient, subject, content, createdBy, createdOn, expireOn, statusId, attempts, priority, reference FROM moved_rows
		)
		UPDATE dashboard
		SET 
			complete = complete + (SELECT COUNT(ID) FROM moved_rows),
			queue = queue - (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND mnth = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
	_, err := repo.db.Conn.Exec(query)
	if err != nil {
		return false, err
	}

	return success, nil
}
