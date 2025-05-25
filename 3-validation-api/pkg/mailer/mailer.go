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

func SendEmail(emailConfig *configs.EmailConfig, msg *SendMsg) {
	e := email.NewEmail()
	e.From = msg.AddressTo
	e.To = []string{msg.AddressTo}
	e.Subject = msg.SubjectMsg
	e.Text = []byte(msg.TextMsg)
	err := e.Send(
		emailConfig.SmtpAddress,
		smtp.PlainAuth(
			"",
			emailConfig.Email,
			emailConfig.Password,
			emailConfig.Address,
		),
	)
	if err != nil {
		log.Fatal("err: ", err)
	}
}
