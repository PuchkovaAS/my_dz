package mailer

import (
	"3-validation-api/configs"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type SendMsg struct {
	AddressTo  string
	SubjectMsg string
	TextMsg    string
}

func SendEmail(emailConfig *configs.EmailConfig, msg *SendMsg) error {
	e := email.NewEmail()
	e.From = emailConfig.Email
	e.To = []string{msg.AddressTo}
	e.Subject = msg.SubjectMsg
	e.Text = []byte(msg.TextMsg)
	err := e.Send(
		emailConfig.SmtpAddress,
		smtp.PlainAuth(
			"",
			emailConfig.Email,
			emailConfig.Password,
			emailConfig.SmtpServer,
		),
	)
	if err != nil {
		log.Printf("err: %s", err)
	}
	return err
}
