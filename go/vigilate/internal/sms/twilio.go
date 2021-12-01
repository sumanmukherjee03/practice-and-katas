package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tsawler/vigilate/internal/config"
)

func SendTextWithTwilio(to string, message string, app *config.AppConfig) error {
	secret := app.PreferenceMap["twilio_auth_token"]
	key := app.PreferenceMap["twilio_sid"]
	urlString := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", key)

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", app.PreferenceMap["twilio_phone_number"])
	msgData.Set("Body", message)

	// msgData.encode() encodes the url values into an url encoded form
	// strings.NewReader is similar to a strings.NewBufferString, except that it is a readonly source
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlString, &msgDataReader)
	if err != nil {
		log.Error(fmt.Errorf("ERROR : Problem generating http request to twilio for outbound sms - %v", err))
		return err
	}

	req.SetBasicAuth(key, secret)
	// Accept a json response
	req.Header.Add("Accept", "application/json")
	// Send POST request using url encoded form data
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Errorf("ERROR : Could not successfully send a request to twilio for outbound sms - %v", err))
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			log.Error(fmt.Errorf("ERROR : Could not parse json response from twilio - %v", err))
			return err
		}
	} else {
		twilioErr := fmt.Errorf("ERROR : Could not send sms because twilio returned an error status code - %d", resp.StatusCode)
		log.Error(twilioErr)
		return twilioErr
	}

	return nil
}
