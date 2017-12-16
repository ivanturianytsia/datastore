package main

import (
	"encoding/json"
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
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	var body map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return err
	}
	if _, ok := body["error"]; ok {
		if v, ok := body["message"].(string); ok {
			return fmt.Errorf(v)
		} else {
			return fmt.Errorf("Error while sending SMS with code")
		}
	}
	return nil
}

func (service *SMSApiService) SendCode(user User, code string) error {
	to := user.PhoneNumber
	message := fmt.Sprintf("%s", code)

	return service.makeRequest(fmt.Sprintf("https://api.smsapi.com/sms.do?to=%s&message=%s&format=json", to, message))
}
