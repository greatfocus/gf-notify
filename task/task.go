package task

import (
	"context"
	"log"
	"time"

	"github.com/greatfocus/gf-notify/repositories"
	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-sframe/server"
)

// Tasks struct
type Tasks struct {
	messageRepository *repositories.MessageRepository
	db                *database.Conn
	Timeout           uint64
}

// Init required parameters
func (t *Tasks) Init(s *server.Meta) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(s.DB, s.Cache)
	t.db = s.DB
	t.Timeout = s.Timeout
}

// SendQueuedEmails intiates the job to send queued messages
func (t *Tasks) SendQueuedEmails() {
	// log.Println("Scheduler_SendQueuedEmails started")
	// request := services.EmailService{}
	// services.SendQueuedEmails(t.messageRepository, &request)
	// log.Println("Scheduler_SendQueuedEmails ended")
}

// MoveStagedToQueue ...
/**
This job is responsible for moving staged messages to the queue
1. staged is a tray containing all messages from the api
2. queue are all messages ready for processing
**/
func (t *Tasks) MoveStagedToQueue() {
	log.Println("Scheduler_MoveStagedToQueue started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	success, err := t.messageRepository.MoveStagedToQueue(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()

	success, err := t.messageRepository.ReQueue(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	success, err := t.messageRepository.MoveFailedQueue(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	success, err := t.messageRepository.MoveCompleteQueue(ctx)
	if err != nil && !success {
		log.Println("Scheduler_MoveOutCompleteQueue failed")
		return
	}
	log.Println("Scheduler_MoveOutCompleteQueue succeded")
}
