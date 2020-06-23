package task

import (
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/repositories"
)

// MessageOut struct
type MessageOut struct {
	messageRepository *repositories.MessageRepository
}

// Start intiates the job
func (m *MessageOut) Start(db *database.DB) {
	m.messageRepository = &repositories.MessageRepository{}
	m.messageRepository.Init(db)
	return
}
