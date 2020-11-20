package email_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/messager/email"
	"github.com/stretchr/testify/assert"
)

func TestEmailSend(t *testing.T) {
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	if host == "" {
		return
	}

	portI, err := strconv.Atoi(port)
	assert.NoError(t, err)

	client := email.NewClient(host, portI, user, password)
	assert.NoError(t, client.Send("Hello, world", "This is message body", "mylxsw@aicode.cc", "mylxsw@126.com"))
}
