package connector_test

import (
	"context"
	"testing"
	"time"

	"github.com/mylxsw/adanos-alert/pkg/connector"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	assert.NoError(t, connector.NewConnector("", "http://localhost:19999").Send(
		ctx,
		connector.NewEvent("Hello, world").
			WithMeta("occur_at", time.Now()).
			WithMeta("user", "adanos").
			WithTags("hello", "connector").
			WithOrigin("connector"),
	))
}
