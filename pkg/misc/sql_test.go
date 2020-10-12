package misc_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/stretchr/testify/assert"
)

func TestSQLFinger(t *testing.T) {
	assert.Equal(t, "xxxxxx adsfa", misc.SQLFinger("XXXXxx adsfa"))
	assert.Equal(t, "select id , name from users where id in ( ... ) and age > ?", misc.SQLFinger("Select id, name from users where id in (1, 2,3 ,4 ) and age > 19"))
}
