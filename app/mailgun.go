package main

import (
	"fmt"

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

	if _, _, err := mg.Send(message); err != nil {
		return err
	}
	return nil
}
