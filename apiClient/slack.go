package apiClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aleph-icabod/json-placeholder/entity"
	"io/ioutil"
	"net/http"
)

// Slack struct to abstract the behavior for slack API
type Slack struct {
	webhookURL string
	client     *http.Client
}

type Attachment struct {
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

func NewSlack(url string) *Slack {
	return &Slack{
		webhookURL: url,
		client:     http.DefaultClient,
	}
}

// SendNotification creates a new notification on slack with webhook
func (s *Slack) SendNotification(msg string, payload *entity.Photo) error {

	message := struct {
		Attachments []Attachment `json:"attachments"`
	}{
		Attachments: []Attachment{
			{
				Title:    msg,
				MrkdwnIn: []string{"text"},
			},
		},
	}

	message.Attachments[0].Text = fmt.Sprintf("```%v```", payload.ToJsonString())

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s", s.webhookURL),
		bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	response, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Println("Slack response: ", string(data))
	return nil
}
