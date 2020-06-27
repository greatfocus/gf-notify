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
		ChannelID: 2,
		StatusID:  1,
		Page:      1,
	}
	log.Println("Fetching Email messages available to send")
	msgs, err := repo.GetMessages("queue", params)
	if err != nil {
		log.Println("Error fetching Email messages available to send")
	} else {
		queueMessages(repo, msgs, request)
		SendBulk(request)
	}
}

// queueMessages update queue messages
func queueMessages(repo *repositories.MessageRepository, msgs []models.Message, request Request) {
	log.Println("Updating queued Email messages available to send")
	var wg sync.WaitGroup
	recipient := make([]string, len(msgs))
	messages := make([]string, len(msgs))
	for i := 0; i < len(msgs); i++ {
		wg.Add(1)
		msgs[i].StatusID = 2
		recipient[i] = msgs[i].Recipient
		messages[i] = msgs[i].Content

		go updateMessage(&wg, repo, "queue", msgs[i])
	}

	request.Recipients = recipient
	request.Messages = messages

	wg.Wait()
}

// updateMessage change message status
func updateMessage(wg *sync.WaitGroup, repo *repositories.MessageRepository, table string, message models.Message) {
	defer wg.Done()
	_, err := repo.Update(table, message)
	if err != nil {
		log.Println("Failed to update Email message with ID", message.ID)
	}
}
