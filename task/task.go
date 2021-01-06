package task

import (
	"log"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-notify/services"
)

// Tasks struct
type Tasks struct {
	messageRepository *repositories.MessageRepository
	config            *config.Config
	db                *database.Conn
}

// Init required parameters
func (t *Tasks) Init(db *database.Conn, cache *cache.Cache, config *config.Config) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(db, cache)
	t.config = config
	t.db = db
}

// SendQueuedEmails intiates the job to send queued messages
func (t *Tasks) SendQueuedEmails() {
	log.Println("Scheduler_SendQueuedEmails started")
	request := services.EmailService{
		Host:     t.config.Integrations.Email.Host,
		Port:     t.config.Integrations.Email.Port,
		From:     t.config.Integrations.Email.From,
		User:     t.config.Integrations.Email.User,
		Password: t.config.Integrations.Email.Password,
	}
	services.SendQueuedEmails(t.messageRepository, &request)
	log.Println("Scheduler_SendQueuedEmails ended")
}

// MoveStagedToQueue ...
/**
This job is responsible for moving staged messages to the queue
1. staged is a tray containing all messages from the api
2. queue are all messages ready for processing
**/
func (t *Tasks) MoveStagedToQueue() {
	log.Println("Scheduler_MoveStagedToQueue started")
	success, err := t.messageRepository.MoveStagedToQueue()
	if err != nil && !success {
		log.Println("Scheduler_MoveStagedToQueue failed")
		return
	}
	log.Println("Scheduler_MoveStagedToQueue succeded")
}

// ReQueueProcessingEmails ...
/**
This job is responsible for re-queueing messages that have been processing for more than 10mins
1. This may occur due to service being stopped
2. or panic error withing golang
**/
func (t *Tasks) ReQueueProcessingEmails() {
	log.Println("Scheduler_ReQueueProcessingEmails started")
	success, err := t.messageRepository.ReQueueProcessingEmails()
	if err != nil && !success {
		log.Println("Scheduler_ReQueueProcessingEmails failed")
		return
	}
	log.Println("Scheduler_ReQueueProcessingEmails succeded")
}

// MoveOutFailedQueue ...
/**
This job is responsible for moving failed messages to the failed
1. queue is a tray containing all messages been sent out
2. failed are all messages that failed to sent with attempt greater >= 5
**/
func (t *Tasks) MoveOutFailedQueue() {
	log.Println("Scheduler_MoveOutFailedQueue started")
	success, err := t.messageRepository.MoveOutFailedQueue()
	if err != nil && !success {
		log.Println("Scheduler_MoveOutFailedQueue failed")
		return
	}
	log.Println("Scheduler_MoveOutFailedQueue succeded")
}

// MoveOutCompleteQueue ...
/**
This job is responsible for moving completed messages to the complete
1. queue is a tray containing all messages been sent out
2. completed are all messages that succeded to sent with attempt less < 5
**/
func (t *Tasks) MoveOutCompleteQueue() {
	log.Println("Scheduler_MoveOutCompleteQueue started")
	success, err := t.messageRepository.MoveOutCompleteQueue()
	if err != nil && !success {
		log.Println("Scheduler_MoveOutCompleteQueue failed")
		return
	}
	log.Println("Scheduler_MoveOutCompleteQueue succeded")
}

/**
Running database script is important for the following
1. To archive archiving of data in the tables, we create tables every month
2. This creates new tables for the new month and reduce the database load to query
3. We have also split the tables into staging, queue, failed and done to avoid database deadlocks
	- staging: new messages from the API go in here. This reduces deadlock in jobs since http bulk messages can cause performance isssues
	- queue: messages are moved here as current messages being sent. This helps to isolate process
	- failed: all failed messages go in here, this helps to wipe and reduce the queue
	- complete: all successful messages are isolated here, this helps with reports isolation
4. This breadown structure also helps with generating a proper dashboard for messages and reporting
**/

// RunDatabaseScripts intiates running database scripts
func (t *Tasks) RunDatabaseScripts() {
	log.Println("Scheduler_RunDatabaseScripts started")
	if t.config.Database.Master.ExecuteSchema {
		t.db.Master.ExecuteSchema(t.db.Master.Conn)
	}
	if t.config.Database.Slave.ExecuteSchema {
		t.db.Slave.ExecuteSchema(t.db.Slave.Conn)
	}
	log.Println("Scheduler_RunDatabaseScripts ended")
}

// RebuildIndexes make changes to indexes
func (t *Tasks) RebuildIndexes() {
	log.Println("Scheduler_RebuildIndexes started")
	if t.config.Database.Master.ExecuteSchema {
		t.db.Master.RebuildIndexes(t.db.Master.Conn, "gf_query")
	}
	if t.config.Database.Slave.ExecuteSchema {
		t.db.Slave.RebuildIndexes(t.db.Slave.Conn, "gf_query")
	}
	log.Println("Scheduler_RebuildIndexes ended")
}
