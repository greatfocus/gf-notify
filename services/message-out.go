package services

import (
	"log"
	"sync"

	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

/**
HERE IS THE LOGIN INVOLVED
1. Theory here is to send every messages per second
2. Updated the messaged status from requested to queue
3. Send email for the messages in bulks of number of seconds
4. Update the messages status incase of any failure or success
**/

// SendNewEmails
func SendNewEmails(repo *repositories.MessageRepository, request Request) {
	params := repositories.MessageParam{
		ChannelID: 1,
		StatusID:  1,
	}
	log.Println("Fetching Email messages available to send")
	msgs, err := repo.GetMessages(params)
	if err != nil {
		log.Println("No Email messages available to send")
	} else {
		queueMessages(repo, msgs, request)
		SendBulk(request)
	}
}

// queueMessages update queue messages
func queueMessages(repo *repositories.MessageRepository, msgs []models.Message, request Request) {
	log.Println("Updating queued Email messages available to send")
	var wg sync.WaitGroup
	for i := 1; i <= len(msgs); i++ {
		wg.Add(1)
		msgs[i].StatusID = 2
		request.Messages[i] = msgs[i].Content
		request.Recipients[i] = msgs[i].Recipient
		go updateMessage(&wg, repo, msgs[i])
	}
	wg.Wait()
}

// updateMessage change message status
func updateMessage(wg *sync.WaitGroup, repo *repositories.MessageRepository, message models.Message) {
	defer wg.Done()
	err := repo.UpdateMessage(message)
	if err != nil {
		log.Println("Failed to update Email message with ID", message.ID)
	}
}
