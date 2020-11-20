package email

import "gopkg.in/gomail.v2"

// Client is a email sender client
type Client struct {
	dailer *gomail.Dialer
	from   string
}

// NewClient create a new mail client
func NewClient(host string, port int, username string, password string) *Client {
	dailer := gomail.NewDialer(host, port, username, password)
	return &Client{dailer: dailer, from: username}
}

// Send send email to users
func (m Client) Send(subject, body string, users ...string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", users...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dailer.DialAndSend(msg)
}
