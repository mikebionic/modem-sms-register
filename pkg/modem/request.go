package modem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type SMSPayload struct {
	PhoneNumber string `json:"phone_number"`
	MessageText string `json:"message_text"`
	Timestamp   string `json:"timestamp"`
}

func Send(urlAddress, phoneNumber, messageText, shaKey string) error {
	payload := SMSPayload{
		PhoneNumber: phoneNumber,
		MessageText: messageText,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", urlAddress, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-access-token", shaKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	log.WithFields(log.Fields{
		"status_code": resp.StatusCode,
	}).Debug("HTTP request successful")
	return nil
}
