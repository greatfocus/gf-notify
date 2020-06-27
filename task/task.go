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
}

// Init required parameters
func (t *Tasks) Init(db *database.DB, config *config.Config) {
	t.messageRepository = &repositories.MessageRepository{}
	t.messageRepository.Init(db)
	t.config = config
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
