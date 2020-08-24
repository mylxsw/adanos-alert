package misc_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/stretchr/testify/assert"
)

func TestIfElse(t *testing.T) {
	assert.Equal(t, "pos", misc.IfElse(true, "pos", "neg"))
	assert.Equal(t, "neg", misc.IfElse(false, "pos", "neg"))
}
