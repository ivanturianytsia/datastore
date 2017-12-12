package main

import (
	"fmt"
	"net/http"
	"os"
)

type SMSApiService struct {
	token string
}

func NewSMSApiService() CodeAuthService {
	token := os.Getenv("SMSAPI_TOKEN")

	return &SMSApiService{
		token: token,
	}
}

func (service *SMSApiService) makeRequest(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", service.token))
	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}

func (service *SMSApiService) SendCode(user User, code string) error {
	to := user.PhoneNumber
	message := fmt.Sprintf("%s", code)

	return service.makeRequest(fmt.Sprintf("https://api.smsapi.com/sms.do?to=%s&message=%s&format=json", to, message))
}
