package matcher_test

import (
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/internal/matcher"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMessageFinger(t *testing.T) {
	var msg = repository.Message{
		ID:      primitive.NewObjectID(),
		Content: `{"log_level": "debug", "message": "request", "context": {"user_id": 123}}`,
		Meta: repository.MessageMeta{
			"environment": "dev",
			"server":      "192.168.1.1",
		},
		Tags:      []string{"php", "nodejs"},
		Origin:    "Filebeat",
		CreatedAt: time.Now(),
	}

	{
		f, err := matcher.NewMessageFinger(`Meta["server"] + ":" + Origin`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "192.168.1.1:Filebeat", finger)
	}

	{

		f, err := matcher.NewMessageFinger(``)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "", finger)
	}

	{
		f, err := matcher.NewMessageFinger(`"hello world"`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", finger)
	}

	{
		f, err := matcher.NewMessageFinger(`Meta["not_exist_key"]`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "", finger)
	}

	{
		f, err := matcher.NewMessageFinger(`124`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "124", finger)
	}

	{
		f, err := matcher.NewMessageFinger(`true`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "true", finger)
	}
}
