package expo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aeramu/menfess-server/internal/notification/service"

	log "github.com/sirupsen/logrus"
)

const (
	baseUrl     = "https://exp.host/--/api/v2/push/send"
	contentType = "application/json"
)

func NewClient() service.PushServiceClient {
	return &client{
		http: http.DefaultClient,
	}
}

type client struct {
	http *http.Client
}

func (c *client) Send(ctx context.Context, pushToken map[string]bool, title string, body string, data string) error {
	var tokens []string
	for key, _ := range pushToken {
		tokens = append(tokens, key)
	}

	reqBody := requestBody{
		To:    tokens,
		Title: title,
		Body:  body,
		Data:  data,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"title": title,
			"body":  body,
			"data":  data,
		}).Errorln("[Send] Failed marshal request body to json")
		return err
	}

	_, err = c.http.Post(baseUrl, contentType, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"title": title,
			"body":  body,
			"data":  data,
		}).Errorln("[Send] Failed send request to handler")
		return err
	}

	return nil

}

type requestBody struct {
	To    []string `json:"to"`
	Title string   `json:"title,omitempty"`
	Body  string   `json:"body"`
	Data  string   `json:"data,omitempty"`
}
