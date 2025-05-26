package mailer

import (
	"3-validation-api/configs"
	"fmt"
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
	if emailConfig.SenderName != "" {
		e.From = fmt.Sprintf(
			"%s <%s>",
			emailConfig.SenderName,
			emailConfig.Email,
		)
	} else {
		e.From = emailConfig.Email
	}
	e.To = []string{msg.AddressTo}
	e.Subject = msg.SubjectMsg
	e.Text = []byte(msg.TextMsg)
	err := e.Send(
		fmt.Sprintf("%s:%s", emailConfig.SmtpHost, emailConfig.SmtpPort),
		smtp.PlainAuth(
			"",
			emailConfig.Email,
			emailConfig.Password,
			emailConfig.SmtpHost,
		),
	)
	if err != nil {
		log.Printf("err: %s", err)
	}
	return err
}
