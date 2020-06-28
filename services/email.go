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
	User       string
	Password   string
	From       string
	Messages   []string
	Recipients []string
	Status     []bool
}

// Send initiates sending of the messages
func Send(i int, email *Request, wg *sync.WaitGroup) bool {
	sent := true
	defer wg.Done()
	// hostname is used by PlainAuth to validate the TLS certificate.
	to := []string{email.Recipients[i]}
	msg := []byte(email.Messages[i])
	auth := smtp.PlainAuth("", email.User, email.Password, email.Host)
	err := smtp.SendMail(email.Host+":"+email.Port, auth, email.From, to, msg)
	if err != nil {
		log.Println(err)
		sent = false
	}
	email.Status[i] = sent
	return sent
}
