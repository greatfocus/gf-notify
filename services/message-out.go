package services

import (
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-notify/repositories"
)

// MessageOut struct
type MessageOut struct {
	notifyRepository *repositories.NotifyRepository
}

// Init method
func (m *MessageOut) Init(db *database.DB) {
	m.notifyRepository = &repositories.NotifyRepository{}
	m.notifyRepository.Init(db)
}

// Start method for starting the message queue
func (m *MessageOut) Start() {
	// start sending messages
}
