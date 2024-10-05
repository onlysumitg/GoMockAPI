package models

import (
	"github.com/onlysumitg/GoMockAPI/env"
	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailRequest struct {
	To       []string
	Subject  string
	Body     string
	Template string
	Data     any
}

type AccountEmailTemplateData struct {
	User            *User
	VerficationId   string
	VerficationLink string
}

// ------------------------------------------------------
//
// ------------------------------------------------------
func SetupMailServer() *mail.SMTPServer {
	host, port, user, pwd := env.GetSMTPServerData()
	server := mail.NewSMTPClient()
	server.Host = host
	server.Port = port // SMTP Port 	465 (25 or 587 for non-SSL)
	server.Username = user
	server.Password = pwd
	server.Encryption = mail.EncryptionTLS
	return server
}
