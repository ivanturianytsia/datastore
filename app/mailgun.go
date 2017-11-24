package main

import (
	"fmt"
	"log"

	mailgun "github.com/mailgun/mailgun-go"
)

var (
	mg mailgun.Mailgun
)

func sendCode(to, code string) error {
	message := mailgun.NewMessage(
		"passwordless@mail.nomidigital.com",
		"AGH Datastore Auth Code",
		fmt.Sprintf("Your auth code is: %s. Note that we never send you links in our emails. If it is not you, let us know.", code),
		to,
	)

	resp, id, err := mg.Send(message)
	if err != nil {
		return err
	}
	log.Printf("Auth code message sent to %s!\nID: %s\nResponse: %s\n", to, id, resp)
	return nil
}
