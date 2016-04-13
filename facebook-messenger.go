package facebookMessenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// FacebookMessenger ...
type FacebookMessenger struct {
	Token string
}

// CallbackMessage ...
type CallbackMessage struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        int          `json:"id"`
	Time      int          `json:"time"`
	Messaging []*Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    *ID       `json:"sender"`
	Recipient *ID       `json:"recipient"`
	Timestamp int       `json:"timestamp"`
	Message   *Message  `json:"message,omitempty"`
	Delivery  *Delivery `json:"delivery,omitempty"`
	string    `json:""`
}

// Message ...
type Message struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// Delivery ...
type Delivery struct {
	Mids      []string `json:"mids"`
	Watermark int      `json:"watermark"`
	Seq       int      `json:"seq"`
}

// ID ...
type ID struct {
	ID int `json:"id"`
}

// Text ...
type Text struct {
	Text string `json:"text"`
}

// SendMessage ...
type SendMessage struct {
	Recipient *ID   `json:"recipient"`
	Message   *Text `json:"message"`
}

// SendTextMessage ...
func (fb *FacebookMessenger) SendTextMessage(recipient int, text string) error {

	m := &SendMessage{
		Recipient: &ID{ID: recipient},
		Message:   &Text{Text: text},
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages?access_token="+fb.Token, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	log.Print("Response: ", result)
	return nil

}
