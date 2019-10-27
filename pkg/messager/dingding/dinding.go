package dingding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DingdingMessage is a message holds all informations for a dingding sender
type DingdingMessage struct {
	Message MarkdownMessage `json:"message"`
	Token   string          `json:"token"`
}

func (dm *DingdingMessage) Encode() []byte {
	data, _ := json.Marshal(dm)
	return data
}

func (dm *DingdingMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &dm)
}

// MarkdownMessage is a markdown message for dingding
type MarkdownMessage struct {
	Type     string              `json:"msgtype,omitempty"`
	Markdown MarkdownMessageBody `json:"markdown,omitempty"`
	At       MessageAtSomebody   `json:"at,omitempty"`
}

// Encode encode markdown message to json bytes
func (m MarkdownMessage) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// NewMarkdownMessage create a new MarkdownMessage
func NewMarkdownMessage(title string, body string, mobiles []string) MarkdownMessage {
	return MarkdownMessage{
		Type: "markdown",
		Markdown: MarkdownMessageBody{
			Title: title,
			Text:  body,
		},
		At: MessageAtSomebody{
			Mobiles: mobiles,
		},
	}
}

// MarkdownMessageBody is markdown body
type MarkdownMessageBody struct {
	Title      string `json:"title,omitempty"`
	Text       string `json:"text,omitempty"`
	MessageURL string `json:"messageUrl,omitempty"`
}

// MessageAtSomebody @ someone
type MessageAtSomebody struct {
	Mobiles []string `json:"atMobiles"`
	ToAll   bool     `json:"isAtAll"`
}

type Dingding struct {
	Endpoint string
	Token    string
}

func NewDingding(token string) *Dingding {
	return &Dingding{Endpoint: "https://oapi.dingtalk.com/robot/send?access_token=", Token: token}
}

type Message interface {
	Encode() ([]byte, error)
}

// dingResponse 钉钉响应
type dingResponse struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

func (ding *Dingding) Send(msg Message) error {
	url := ding.Endpoint + ding.Token

	msgEncoded, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("dingding message encode failed: %s", err.Error())
	}

	reader := bytes.NewReader(msgEncoded)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return fmt.Errorf("dingding create request failed: %w", err)
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("dingding send msg failed: %w", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("dingding read response failed: %w", err)
	}

	var dresp dingResponse
	if err := json.Unmarshal(respBytes, &dresp); err != nil {
		return fmt.Errorf("send finished, response： %s", string(respBytes))
	}

	if dresp.ErrorCode > 0 {
		return fmt.Errorf("[%d] %s", dresp.ErrorCode, dresp.ErrorMessage)
	}

	return nil
}
