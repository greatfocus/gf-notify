package task

import (
	"log"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/repositories"
)

// Tasks struct
type Tasks struct {
	messageRepository *repositories.MessageRepository
}

// Init required parameters
func (t *Tasks) Init(db *database.DB) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(db)
}

// Start intiates the job
func (t *Tasks) Start() {
	log.Println("CScheduler started")
	return
}
