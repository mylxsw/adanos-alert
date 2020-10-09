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
	var msg = repository.Event{
		ID:      primitive.NewObjectID(),
		Content: `{"log_level": "debug", "message": "request", "context": {"user_id": 123}}`,
		Meta: repository.EventMeta{
			"environment": "dev",
			"server":      "192.168.1.1",
		},
		Tags:      []string{"php", "nodejs"},
		Origin:    "Filebeat",
		CreatedAt: time.Now(),
	}

	{
		f, err := matcher.NewEventFinger(`Meta["server"] + ":" + Origin`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "192.168.1.1:Filebeat", finger)
	}

	{

		f, err := matcher.NewEventFinger(``)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "", finger)
	}

	{
		f, err := matcher.NewEventFinger(`"hello world"`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", finger)
	}

	{
		f, err := matcher.NewEventFinger(`Meta["not_exist_key"]`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "", finger)
	}

	{
		f, err := matcher.NewEventFinger(`124`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "124", finger)
	}

	{
		f, err := matcher.NewEventFinger(`true`)
		assert.NoError(t, err)

		finger, err := f.Run(msg)
		assert.NoError(t, err)
		assert.Equal(t, "true", finger)
	}
}
