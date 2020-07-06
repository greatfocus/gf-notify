package services

import (
	"log"
	"sync"

	"github.com/greatfocus/gf-notify/models"
	"github.com/greatfocus/gf-notify/repositories"
)

// SendQueuedEmails ...
/**
HERE IS THE LOGIN INVOLVED
1. Theory here is to send every messages per second
2. Updated the messaged status from requested to queue
3. Send email for the messages in bulks of number of seconds
4. Update the messages status incase of any failure or success
**/
func SendQueuedEmails(repo *repositories.MessageRepository, request *EmailRequest) {
	params := repositories.MessageParam{
		ChannelID: 2,
		StatusID:  2,
		Attempts:  5,
		Page:      1,
	}
	log.Println("Scheduler_SendQueuedEmails Fetching Email queued messages")
	msgs, err := repo.GetMessages("queue", params)
	if err != nil {
		log.Println("Scheduler_SendQueuedEmails Error fetching Email queued")
		return
	}

	if len(msgs) > 0 {
		prepareQueueMessages(repo, msgs, request)
		sendBulkEmails(repo, msgs, request)
	} else {
		log.Println("Scheduler_SendQueuedEmails Email queued is empty")
	}
}

// prepareQueueMessages creates post model
func prepareQueueMessages(repo *repositories.MessageRepository, msgs []models.Message, request *EmailRequest) {
	var args []interface{}
	recipient := make([]string, len(msgs))
	subjects := make([]string, len(msgs))
	messages := make([]string, len(msgs))
	status := make([]bool, len(msgs))
	for i := 0; i < len(msgs); i++ {
		recipient[i] = msgs[i].Recipient
		messages[i] = msgs[i].Content
		subjects[i] = msgs[i].Subject
		args = append(args, msgs[i].ID)
	}
	request.Recipients = recipient
	request.Subjects = subjects
	request.Messages = messages
	request.Status = status

	repo.UpdateQueueToProcessing("queue", args)
}

// SendBulk initiates sending of the messages
func sendBulkEmails(repo *repositories.MessageRepository, msgs []models.Message, email *EmailRequest) {
	log.Println("Scheduler_SendQueuedEmails Sending bulk Email messages")
	var wg sync.WaitGroup

	for i := 0; i < len(email.Recipients); i++ {
		wg.Add(1)
		go SendMail(i, email, &wg)
	}

	wg.Wait()
	updateQueueEmail(repo, msgs, email)
}

// updateMessage change message status
func updateQueueEmail(repo *repositories.MessageRepository, msgs []models.Message, email *EmailRequest) {
	for i := 0; i < len(email.Recipients); i++ {
		// check status of email sent
		trial := msgs[i].Attempts + 1
		msgs[i].Attempts = trial
		if email.Status[i] {
			msgs[i].StatusID = 4
			msgs[i].Reference = "1"
		} else {
			msgs[i].Reference = "0"
			if msgs[i].Attempts < 5 {
				msgs[i].StatusID = 2
			} else {
				msgs[i].StatusID = 5
			}
		}
		_, err := repo.Update("queue", msgs[i])
		if err != nil {
			log.Println("Failed to update Email message with ID", msgs[i].ID)
			log.Println(err)
		}
	}
}
