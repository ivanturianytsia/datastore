package main

import (
	"fmt"
	"os"

	mailgun "github.com/mailgun/mailgun-go"
)

type CodeAuthService interface {
	SendCode(to, code string) error
}

type MailgunService struct {
	mg mailgun.Mailgun
}

func NewMailgunService() CodeAuthService {
	mg := mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_KEY"),
		os.Getenv("MAILGUN_PUBKEY"))

	return &MailgunService{
		mg: mg,
	}
}

func (service *MailgunService) SendCode(to, code string) error {
	message := mailgun.NewMessage(
		"passwordless@mail.nomidigital.com",
		"AGH Datastore Auth Code",
		fmt.Sprintf("Your auth code is: %s. Note that we never send you links in our emails. If it is not you, let us know.", code),
		to,
	)

	if _, _, err := service.mg.Send(message); err != nil {
		return err
	}
	return nil
}
