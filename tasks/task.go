package tasks

import (
	"github.com/greatfocus/gf-notify/services"
	"github.com/greatfocus/go-frame/database"
)

// MessageOut queues messages to start sending
func MessageOut(db *database.DB) {
	// Initialize controller
	messageOut := services.MessageOut{}
	messageOut.Init(db)
	messageOut.Start()
}
