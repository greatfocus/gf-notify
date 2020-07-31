package services

import (
	"log"
	"net/smtp"
	"sync"
)

// EmailRequest struct
type EmailRequest struct {
	Host       string
	Port       string
	User       string
	Password   string
	From       string
	Subjects   []string
	Messages   []string
	Recipients []string
	Status     []bool
}

// Auth get the authentication for each batch to avoid having too many logins
func Auth(email *EmailRequest) smtp.Auth {
	auth := smtp.PlainAuth("", email.User, email.Password, email.Host)
	return auth
}

// SendMail initiates sending of the email
func SendMail(i int, email *EmailRequest, auth smtp.Auth, wg *sync.WaitGroup) {
	sent := true
	// hostname is used by PlainAuth to validate the TLS certificate.
	to := []string{email.Recipients[i]}
	msg := []byte("To: " + email.Recipients[i] + "\r\n" +
		"Subject: " + email.Subjects[i] + "\r\n" +
		"\r\n" +
		email.Messages[i] + ".\r\n")

	// Please watch out here not to have golang panic error. This means there might be some jobs in progres that may never resolve
	err := smtp.SendMail(email.Host+":"+email.Port, auth, email.From, to, msg)
	if err != nil {
		log.Println(err)
		sent = false
	}
	email.Status[i] = sent
	wg.Done()
}
