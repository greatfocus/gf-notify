package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/greatfocus/gf-notify/models"
	cache "github.com/greatfocus/gf-sframe/cache"
	"github.com/greatfocus/gf-sframe/database"
)

// MessageRepository struct
type MessageRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// QueryParam struct
type QueryParam struct {
	Table     string
	Month     string
	Year      string
	Statement string
	Args      []interface{}
	Repo      *MessageRepository
	ctx       context.Context
	enKey     string
}

// MessageParam struct
type MessageParam struct {
	ChannelID string
	Status    string
	Attempts  int64
	LastID    string
}

// Init method
func (repo *MessageRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
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

/*
*
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
*
*/
func getTable(name string) string {
	year, month := getYearAndMonth()
	tableName := name + year + month
	return tableName
}

/*
*
There are several tables involved in the logic for messaging
We use this generic function to concatinate the table and date
*
*/
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
		err := rows.Scan(&message.ID, &message.ChannelID, &message.Recipient, &message.Subject, &message.Content,
			&message.CreatedOn, &message.ExpireOn, &message.Status, &message.Attempts, &message.Priority)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// insert adds new  database record
func insert(params *QueryParam) error {
	_, inserted := params.Repo.db.Insert(params.ctx, params.Statement, params.Args...)
	if !inserted {
		return errors.New("create failed")
	}
	return nil
}

// update makes changes to database record
func update(params *QueryParam) error {
	updated := params.Repo.db.Update(params.ctx, params.Statement, params.Args...)
	if !updated {
		return errors.New("update failed")
	}
	return nil
}

// queryMessages executes to get messages
func query(params *QueryParam) ([]models.Message, error) {
	rows, err := params.Repo.db.Query(params.ctx, params.Statement, params.Args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	return messageMapper(rows)
}

// newQueryParams returns params
func newQueryParams(table string, statement string, args []interface{}, repo *MessageRepository, ctx context.Context, enKey string) *QueryParam {
	// get current month and date
	year, month := getYearAndMonth()
	return &QueryParam{
		Table:     table,
		Month:     month,
		Year:      year,
		Statement: statement,
		Args:      args,
		Repo:      repo,
		ctx:       ctx,
		enKey:     enKey,
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

// ReportMessages method returns messages from the database
func (repo *MessageRepository) Report(ctx context.Context, enKey string, table string, channel string, year string, month string, lastID string) ([]models.Message, error) {
	// prepare the statement
	statement := `
	select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdOn, m.expireOn, m.status, m.attempts, m.priority 
	from $table m
	inner join channel c on c.id = m.channelId
	where m.channelId = $1 and m.id >= $2
	order BY m.createdOn DESC limit 20
	`
	// prepare the args
	var args []interface{}
	args = append(args, channel, lastID)
	tt := year + month
	newTable := fmt.Sprintf("%s%s", table, tt)
	qry := strings.Replace(statement, "$table", string(newTable), -1)
	queryParams := newQueryParams(table, qry, args, repo, ctx, enKey)
	return query(queryParams)
}

// GetMessages method returns messages from the database
func (repo *MessageRepository) GetMessages(ctx context.Context, enKey string, table string, messageParam MessageParam) ([]models.Message, error) {
	// prepare the statement
	var statement string
	var args []interface{}
	if messageParam.LastID != "" {
		statement = `
		select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdOn, m.expireOn, m.status, m.attempts, m.priority 
		from $table m
		inner join channel c on c.id = m.channelId
		where channelId = $1 and m.status=$2 and m.attempts<$3
		order BY m.createdOn DESC limit 20
		`
		args = append(args, messageParam.ChannelID, messageParam.Status, messageParam.Attempts)
	} else {
		statement = `
		select m.id, m.channelId, c.name as channel, m.recipient, m.subject, m.content, m.createdOn, m.expireOn, m.status, m.attempts, m.priority 
		from $table m
		inner join channel c on c.id = m.channelId
		where channelId = $1 and m.status=$2 and m.attempts<$3 and m.id >=$4
		order BY m.createdOn DESC limit 20
		`
		args = append(args, messageParam.ChannelID, messageParam.Status, messageParam.Attempts, messageParam.LastID)
	}
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo, ctx, enKey)
	return query(queryParams)
}

// Create method created new message request
func (repo *MessageRepository) Create(ctx context.Context, enKey string, table string, message models.Message) (models.Message, error) {
	var id = uuid.New().String()
	statement := `
    insert into $table (id, channelId, recipient, subject, content, createdOn, expireOn, status, priority)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
	`

	var args []interface{}
	args = append(args, id, message.ChannelID, message.Recipient, message.Subject, message.Content,
		message.CreatedOn, message.ExpireOn, message.Status, message.Priority)
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo, ctx, enKey)
	err := insert(queryParams)
	if err != nil {
		return message, err
	}

	message.ID = id
	return message, nil
}

// Update method make changes to message
func (repo *MessageRepository) Update(ctx context.Context, enKey string, table string, message models.Message) error {
	statement := `
    update $table
	set 
		status=$2, 
		attempts=$3,
		reference=$4	
    where id=$1
	`

	var args []interface{}
	args = append(args, message.ID, message.Status, message.Attempts, message.Reference)
	queryParams := newQueryParams(table, replaceTableHolder(table, statement), args, repo, ctx, enKey)
	err := update(queryParams)
	if err != nil {
		return err
	}

	return nil
}

// UpdateQueueToProgress method make changes to message
func (repo *MessageRepository) UpdateQueueToProgress(ctx context.Context, table string, args []interface{}) error {
	statement := `
    update queue$tt
	set 
		status="pending"
    where id IN $IN
	`
	newArgs := prepareInStatement(args)
	query := strings.Replace(statement, "$IN", newArgs, -1)
	qry := replaceTimeHolder(query)
	queryParams := newQueryParams(table, replaceTableHolder(table, qry), args, repo, ctx, "")
	err := update(queryParams)
	if err != nil {
		return err
	}

	return nil
}

// ReQueue runs a script to move data
func (repo *MessageRepository) ReQueue(ctx context.Context) (bool, error) {
	statement := `
		UPDATE queue$tt
		SET status="new"
		WHERE status="pending" and EXTRACT(MINUTE FROM updatedOn) > 10;
	`
	query := replaceTimeHolder(statement)
	updated := repo.db.Update(ctx, query)
	if !updated {
		derr := errors.New("message update failed")
		return updated, derr
	}

	return updated, nil
}

// MoveStagedToQueue runs a script to move data
func (repo *MessageRepository) MoveStagedToQueue(ctx context.Context) (bool, error) {
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
				ORDER BY id
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO queue$tt (id, channelId, recipient, subject, content, createdOn, expireOn, status, attempts, priority)
			SELECT id, channelId, recipient, subject, content, createdOn, expireOn, "new", 0, priority FROM moved_rows
		)		
		UPDATE dashboard
		SET 
			queue = queue + (SELECT COUNT(ID) FROM moved_rows), 
			staging = staging - (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND month = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
	updated := repo.db.Update(ctx, query)
	if !updated {
		derr := errors.New("message update failed")
		return updated, derr
	}

	return updated, nil
}

// MoveFailedQueue runs a script to move data
func (repo *MessageRepository) MoveFailedQueue(ctx context.Context) (bool, error) {
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
				WHERE status=5
				ORDER BY id
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO failed$tt (id, channelId, recipient, subject, content, createdOn, expireOn, status, attempts, priority, reference)
			SELECT id, channelId, recipient, subject, content, createdOn, expireOn, status, attempts, priority, reference FROM moved_rows
		)
		UPDATE dashboard
		SET 
			failed = failed + (SELECT COUNT(ID) FROM moved_rows),
			queue = queue - (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND month = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
	updated := repo.db.Update(ctx, query)
	if !updated {
		derr := errors.New("message update failed")
		return updated, derr
	}

	return updated, nil
}

// MoveCompleteQueue runs a script to move data
func (repo *MessageRepository) MoveCompleteQueue(ctx context.Context) (bool, error) {
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
				WHERE status=4
				ORDER BY id
				LIMIT 500)
			RETURNING *
		), insert_moved_rows AS (
			INSERT INTO complete$tt (id, channelId, recipient, subject, content, createdOn, expireOn, status, attempts, priority, reference)
			SELECT id, channelId, recipient, subject, content, createdOn, expireOn, status, attempts, priority, reference FROM moved_rows
		)
		UPDATE dashboard
		SET 
			complete = complete + (SELECT COUNT(ID) FROM moved_rows),
			queue = queue - (SELECT COUNT(ID) FROM moved_rows)
		WHERE year = yr AND month = mnth;
	END $$;
	`

	query := replaceTimeHolder(statement)
	updated := repo.db.Update(ctx, query)
	if updated {
		derr := errors.New("message update failed")
		return updated, derr
	}

	return updated, nil
}
