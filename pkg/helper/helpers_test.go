package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpers_HumanDuration(t *testing.T) {
	assert.Equal(t, "4 天 20 小时 40 分钟", NewHelpers().HumanDuration("7000m"))
	assert.Equal(t, "1 小时 5 分钟", NewHelpers().HumanDuration("65m"))
}
