package dingding_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/template"
	"github.com/mylxsw/adanos-alert/pkg/messager/dingding"
	"github.com/mylxsw/asteria/log"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/stretchr/testify/assert"
)

func TestNewMarkdownMessage(t *testing.T) {
	body := `## Hello

Hello, world！@18344822222 @19999999991 @18888888888
That all.
`
	msg := dingding.NewMarkdownMessage("Hello", body, []string{"18888888888", "12455552211"})
	assert.Equal(t, 4, len(msg.At.Mobiles))
}

func TestDingding_Send(t *testing.T) {

	sql := `
select
      user.id, user.name,
	user.age
FROM users user WHERE id in (1, 2, 3,
5, 6) and enteprise_id = 145;
`

	body := `## <font color="#ea2426">【重要】</font> 这是报警标题 
---
**192.168.1.1:user**

<b><font color="#ea2426">这些查询存在严重的性能问题，影响其它业务服务，请尽快处理!!！</font></b>

- 数据库 **test**，用户 **test**，查询耗时：<font color="#ea2426">10.432 s</font>，解析行数 <font color="#ea2426">1010101</font>，返回行数 **29**

{{sql}}
---
- 数据库 **test**，用户 **test**，查询耗时：<font color="#ea2426">**10.432 s**</font>，解析行数 <font color="#ea2426">**1010101**</font>，返回行数 **29**

{{sql}}
`

	body = strings.ReplaceAll(body, "{{sql}}", "```sql\n"+template.RemoveEmptyLine(sql)+"```\n")

	token := os.Getenv("DINGDING_TOKEN")
	secret := os.Getenv("DINGDING_SECRET")

	fmt.Println(body)

	if token == "" {
		log.Warningf("dingding env is not set")
		return
	}

	message := dingding.NewMarkdownMessage("测试报警标题", body, []string{})
	ding := dingding.NewDingding(token, secret)
	if err := ding.Send(message); err != nil {
		t.Errorf("send failed: %v", err)
	}
}

func TestExtractAtSomeones(t *testing.T) {
	lines := "Hello\n@18888888888\nok,@19999999949\t@19433333334"
	assert.Equal(t, 3, len(dingding.ExtractAtSomeones(lines)))
}
