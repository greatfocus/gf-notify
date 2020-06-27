package services

import (
	"log"
	"net/smtp"
	"sync"
)

// Request struct
type Request struct {
	Host       string
	Port       string
	From       string
	Messages   []string
	Recipients []string
}

// SendBulk initiates sending of the messages
func SendBulk(email Request) {
	log.Println("Sending bulk Email messages available to send")
	var wg sync.WaitGroup

	for i := 1; i <= len(email.Recipients); i++ {
		wg.Add(1)
		go Send(i, email, &wg)
	}

	wg.Wait()
}

// Send initiates sending of the messages
func Send(i int, email Request, wg *sync.WaitGroup) bool {
	defer wg.Done()
	sent := true
	// hostname is used by PlainAuth to validate the TLS certificate.
	to := []string{email.Recipients[i]}
	msg := []byte(email.Messages[i])
	auth := smtp.PlainAuth("", "user@example.com", "password", email.Host)
	err := smtp.SendMail(email.Host+":"+email.Port, auth, email.From, to, msg)
	if err != nil {
		log.Fatal(err)
		sent = false
	}
	return sent
}
