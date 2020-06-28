package task

import (
	"log"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-notify/services"
)

// Tasks struct
type Tasks struct {
	messageRepository *repositories.MessageRepository
	config            *config.Config
	db                *database.DB
}

// Init required parameters
func (t *Tasks) Init(db *database.DB, config *config.Config) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(db)
	t.config = config
	t.db = db
}

// SendNewEmails intiates the job to send new messages
func (t *Tasks) SendNewEmails() {
	log.Println("Scheduler started for new Email Messages")
	request := services.Request{
		Host: t.config.Email.Host,
		Port: t.config.Email.Port,
		From: t.config.Email.From,
	}
	services.SendNewEmails(t.messageRepository, request)
	log.Println("Scheduler stopped for new Email Messages")
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
	var db = database.DB{}
	db.Connect(t.config)
	log.Println("Scheduler_RunDatabaseScripts ended")
}

// MoveStagedToQueue ...
/**
This job is responsible for moving staged messages to the queue
1. staged is a tray containing all messages from the api
2. queue are all messages ready for processing
**/
func (t *Tasks) MoveStagedToQueue() {
	log.Println("Scheduler_MoveStagedToQueue initated")
	success, err := t.messageRepository.MoveStagedToQueue()
	if err != nil && !success {
		log.Println("Scheduler_RunDatabaseScripts failed")
		return
	}
	log.Println("Scheduler_RunDatabaseScripts succeded")
}
