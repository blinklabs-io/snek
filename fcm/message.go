package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/blinklabs-io/snek/internal/logging"
)

type Message struct {
	MessageContent `json:"message"`
}

type MessageContent struct {
	Token        string                 `json:"token"`
	Notification *NotificationContent   `json:"notification,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

type NotificationContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type MessageOption func(*MessageContent)

func WithData(data map[string]interface{}) MessageOption {
	return func(m *MessageContent) {
		m.Data = data
	}
}

func WithNotification(title string, body string) MessageOption {
	return func(m *MessageContent) {
		m.Notification = &NotificationContent{
			Title: title,
			Body:  body,
		}
	}
}

func NewMessage(token string, opts ...MessageOption) *Message {
	if token == "" {
		logging.GetLogger().Fatalf("Token is mandatory for FCM message")
	}

	msg := &Message{
		MessageContent: MessageContent{
			Token: token,
		},
	}
	for _, opt := range opts {
		opt(&msg.MessageContent)
	}
	return msg
}

func Send(accessToken string, projectId string, msg *Message) error {

	fcmEndpoint := fmt.Sprintf(
		"https://fcm.googleapis.com/v1/projects/%s/messages:send",
		projectId,
	)

	// Convert the message to JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(payload))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", fcmEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}
