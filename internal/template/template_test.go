package template

import (
	"fmt"
	"strings"
	"testing"

	pkgJSON "github.com/mylxsw/adanos-alert/pkg/json"
	"github.com/mylxsw/go-ioc"
	"github.com/stretchr/testify/assert"
)

var content = `
create mode 100644 .gitignore
create mode 100644 Gopkg.lock
create mode 100644 Gopkg.toml
create mode 100644 LICENSE
create mode 100644 README.md
create mode 100644 alert/dingding/dingding.go
create mode 100644 api/alert_rule.go
create mode 100644 api/api.go
create mode 100644 api/message.go
create mode 100644 message-view.tmpl
create mode 100644 parser/parser.go
create mode 100644 parser/parser_test.go
create mode 100644 route/api.go
create mode 100644 route/router.go
create mode 100644 route/web.go
create mode 100644 server/main.go
create mode 100644 storage/alert_rules.go
create mode 100644 storage/message.go
create mode 100644 storage/message_group.go
create mode 100644 storage/model.go
create mode 100644 storage/sequence.go
create mode 100644 storage/user.go
create mode 100644 web/message.go`

type testParseData struct {
	Name    string
	Age     int
	Content string
}

func (data testParseData) Strings(ele string) []string {
	return []string{fmt.Sprintf("name=%s", data.Name), fmt.Sprintf("age=%d", data.Age), ele}
}

func TestParse(t *testing.T) {
	{
		parsed, err := Parse(ioc.New(), `{{ range $i, $msg := (explode .Name " ") }} - {{ $msg }} {{ end }}`, testParseData{Name: "普罗米 修斯", Age: 1088})
		if err != nil {
			t.Errorf("test parse template failed: %v", err)
		}

		if parsed != " - 普罗米  - 修斯 " {
			t.Error("test failed")
		}
	}

	{
		parsed, err := Parse(ioc.New(), `{{ range $i, $msg := .Strings "last element" }} - {{ $msg }} {{ end }}`, testParseData{Name: "普罗米 修斯", Age: 1088})
		if err != nil {
			t.Errorf("test parse template failed: %v", err)
		}

		fmt.Println(parsed)
	}

	{
		parsed, err := Parse(ioc.New(), `{{ helpers.JSON .Content "deps.#.first" }}`, testParseData{Name: "Prometheus", Age: 1088, Content: `{"id": 123, "name": "李逍遥", "deps": [{"first": 123, "second": 434}]}`})
		assert.NoError(t, err)
		assert.Equal(t, "[123]", parsed)
	}
}

func TestLeftIdent(t *testing.T) {
	for _, line := range strings.Split(leftIdent("....", content), "\n") {
		if !strings.HasPrefix(line, "....") {
			t.Errorf("leftIdent函数执行异常")
		}
	}
}

var jsonContent = `{
    "message": "sms_send_failed",
    "context": {
        "msg": "短信发送失败，该错误不允许重试其它通道，请人工介入处理",
        "sms": {
            "id": 30139,
            "phone": "15923356841",
            "app_id": 1,
            "template_params": {
                "material_no": "JN027",
                "operator_name": "陈勇",
                "sence_name": "设备维护",
                "node_name": "故障报修",
                "material_name": "机械压机",
                "operate_time": "2018-11-01 23:12:41"
            },
            "template_id": 5,
            "status": 2,
            "created_at": "2018-11-01 23:12:41",
            "updated_at": "2018-11-03 23:31:01"
        },
        "final_channel": "亿美软通",
        "ack": {
            "msg": "停机或空号",
            "code": "MI:0013"
        },
        "line": "/data/webroot/message/app/Components/SmsHandler.php:383"
    },
    "level": 400,
    "level_name": "ERROR",
    "channel": "custom_cmd",
    "datetime": "2018-11-03 23:31:01",
    "extra": {
        "ref": "20181103233101-5bddbf35de697"
    }
}`

func TestJsonGet(t *testing.T) {

	if pkgJSON.Get("message", "", jsonContent) != "sms_send_failed" {
		t.Errorf("test failed")
	}

	if pkgJSON.Get("context.msg", "", jsonContent) != "短信发送失败，该错误不允许重试其它通道，请人工介入处理" {
		t.Errorf("test failed")
	}

}

func TestJsonGets(t *testing.T) {
	var jsonContent1 = `{
	"msg": "abc",
	"context": {
		"c1": { "user": 123 },
		"c2": { "user": 345 }
	}
}`
	if pkgJSON.Gets("msg,message", "", jsonContent1) != "abc" {
		t.Errorf("test failed")
	}

	var jsonContent2 = `{
    "message": "sms_send_failed"
}`
	if pkgJSON.Gets("msg,message", "", jsonContent2) != "sms_send_failed" {
		t.Errorf("test failed")
	}

}

func TestStartsWith(t *testing.T) {
	if !startsWith("stacktrace", "stack") {
		t.Errorf("test failed")
	}

	if !startsWith("stacktrace", "trace", "stack") {
		t.Errorf("test failed")
	}

	if startsWith("stackstrace", "trace") {
		t.Errorf("test failed")
	}
}

func TestEndsWith(t *testing.T) {
	if !endsWith("stacktrace", "trace") {
		t.Errorf("test failed")
	}

	if !endsWith("stacktrace", "stack", "trace") {
		t.Errorf("test failed")
	}

	if endsWith("stackstrace", "test") {
		t.Errorf("test failed")
	}
}

func TestStringMask(t *testing.T) {
	testCases := map[string]string{
		"abcdefg":                  "*******",
		"5e8af5bb09d64979185635bf": "5e8af5************5635bf",
		"4dcd69cf0c98b6d4357e77b7150e56b32ab13afe461839e8934527db36b21091": "4dcd69****************************************************b21091",
	}

	for k, v := range testCases {
		masked := StringMask(k, 6)
		if masked != v {
			t.Errorf("test failed")
		}
	}
}

func TestStringTags(t *testing.T) {
	var testCase = map[string]int{
		"a,b,c,": 3,
		"a,,b,c": 3,
		" ":      0,
		"":       0,
	}

	for k, v := range testCase {
		if len(StringTags(k, ",")) != v {
			t.Error("test failed")
		}
	}
}

func TestRemoveEmptyLine(t *testing.T) {
	var original = `Hello, world
Are you ready?

Nice!          
        
What are you doing?`

	result := RemoveEmptyLine(original)
	fmt.Println(result)
}

func TestMetaFilter(t *testing.T) {
	meta := make(map[string]interface{})
	meta["message.k1"] = "v1"
	meta["message.k2"] = "v2"
	meta["message.k3"] = "v3"
	meta["message.k4"] = "v4"
	meta["message"] = "Hello, world"
	meta["version"] = "1.0"

	var temp = `{{ range $i, $msg := meta_prefix_filter .Meta "message." "version" }}[{{ $i }}: {{ $msg }}]{{ end }}`

	res, err := Parse(ioc.New(), temp, map[string]interface{}{
		"Meta": meta,
	})
	assert.NoError(t, err)
	assert.Equal(t, "[message.k1: v1][message.k2: v2][message.k3: v3][message.k4: v4][version: 1.0]", res)
}

func TestSortMapByKeyHuman(t *testing.T) {
	data := map[string]interface{}{
		"@timestamp":      "123456",
		"a1":              "a1",
		"b1":              "b1",
		"message.message": "hello, world",
		"message.name":    "your name",
		"yyy":             "zzz",
	}

	dataSorted := SortMapByKeyHuman(data)
	for _, d := range dataSorted {
		fmt.Printf("%s: %s\n", d.Key, d.Value)
	}

	assert.Equal(t, len(data), len(dataSorted))
	assert.True(t, strings.HasPrefix(dataSorted[0].Key, "message"))
}

func TestImplode(t *testing.T) {
	assert.Equal(t, "aaa, bbb, ccc", Implode([]string{"aaa", "bbb", "ccc"}, ", "))
	assert.Equal(t, "aaa, bbb, ccc", Implode([]interface{}{"aaa", "bbb", "ccc"}, ", "))
	assert.Equal(t, "aaa, 123, ccc", Implode([]interface{}{"aaa", 123, "ccc"}, ", "))
	assert.Equal(t, "111, 123, 333", Implode([]int{111, 123, 333}, ", "))
	assert.Equal(t, "1322", Implode(1322, ", "))

	type user struct {
		Name string
		Age  int
	}

	users := []user{
		{Name: "zhangsan", Age: 11},
		{Name: "lisi", Age: 22},
	}

	assert.Equal(t, "{zhangsan 11},{lisi 22}", Implode(users, ","))
}

func TestNumberBeauty(t *testing.T) {
	parsed, _ := Parse(ioc.New(), `{{ index .Data "number" | format "%.0f" | number_beauty }} | {{ index .Data "number" | number_beauty }}`, struct {
		Data map[string]interface{}
	}{
		Data: map[string]interface{}{
			"number": 111133448958232.0,
		},
	})
	assert.Equal(t, "111,133,448,958,232 | 111,133,448,958,232.00", parsed)
}

func TestJSONCutOffFields(t *testing.T) {
	data := Serialize(TrimPrefixMapK("context.", MetaFilterPrefix(JSONCutOffFields(20, jsonContent), "context.msg", "context.final_channel")))
	assert.Equal(t, `{"final_channel":"亿美软通","msg":"短信发送失败，该错误不允许重试其它通道，..."}`, data)
}

func TestCutoffLine(t *testing.T) {
	assert.Equal(t, 5+1, len(strings.Split(CutOffLine(5, jsonContent), "\n")))
}

func TestLineFilterInclude(t *testing.T) {
	assert.Equal(t, `        "msg": "短信发送失败，该错误不允许重试其它通道，请人工介入处理",
        "sms": {
            "phone": "15923356841",
            "msg": "停机或空号",`, LineFilterInclude(`"(sms|phone|msg)"`, jsonContent))

	assert.Equal(t, 32, len(strings.Split(LineFilterExclude(`"(sms|phone|msg)"`, jsonContent), "\n")))
}

func TestMarkdown2html(t *testing.T) {
	mc := `# Hello, world

| a | b |
| --- | --- |
| 123 | 456 |

<script>alert('xss')</script>
`
	expected := `<h1>
  Hello, world
</h1>
<table>
  <thead>
    <tr>
      <th>
        a
      </th>
      <th>
        b
      </th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>
        123
      </td>
      <td>
        456
      </td>
    </tr>
  </tbody>
</table>`
	assert.Equal(t, expected, FormatHTML(Markdown2html(mc)))
	assert.Equal(t, "# Hello, world\n\n| a | b |\n| --- | --- |\n| 123 | 456 |", HTML2Markdown(expected))
}

func TestDOMQueryHTMLFirst(t *testing.T) {
	htmlContent := `<html>
 <head></head>
 <body>
  <h2 style="color:#FF0000"> Execution '30' of flow 'ppp' of project 'demo' has failed on Test</h2>
  <table>
   <tbody>
    <tr>
     <td>Start Time</td>
     <td>2020/10/21 22:30:00 CST</td>
    </tr>
    <tr>
     <td>End Time</td>
     <td>2020/10/21 22:30:00 CST</td>
    </tr>
    <tr>
     <td>Duration</td>
     <td>0 sec</td>
    </tr>
    <tr>
     <td>Status</td>
     <td>FAILED</td>
    </tr>
   </tbody>
  </table>
  <a href="http://localhost:8081/executor?execid=30">ppp Execution Link</a>
  <h3>Reason</h3>
  <ul>
   <li><a href="http://localhost:8081/executor?execid=30&job=ppp">Failed job 'ppp' Link</a></li>
   <li>Not running on the assigned executor (any more)</li>
  </ul>
  <h3>Executions from past 72 hours (26 out 26) failed</h3>
  <table>
   <tbody>
    <tr>
     <td>Execution Id</td>
     <td>30</td>
    </tr>
    <tr>
     <td>Start Time</td>
     <td>2020/10/21 22:30:00 CST</td>
    </tr>
    <tr>
     <td>End Time</td>
     <td>2020/10/21 22:30:00 CST</td>
    </tr>
    <tr>
     <td>Status</td>
     <td>FAILED</td>
    </tr>
   </tbody>
  </table>
 </body>
</html>`

	assert.Equal(t, `<li>
  <a href="http://localhost:8081/executor?execid=30&job=ppp">
    Failed job 'ppp' Link
  </a>
</li>
<li>
  Not running on the assigned executor (any more)
</li>`, FormatHTML(DOMFilterHTMLIndex("ul", 0, htmlContent)))
	assert.Equal(t, "Executions from past 72 hours (26 out 26) failed", DOMFilterHTMLIndex("h3", 1, htmlContent))
}

func TestStrConcat(t *testing.T) {
	assert.Equal(t, "s1s2s3", StrConcat("s1", "s2", "s3"))
}

func TestMarkdown2Confluence(t *testing.T) {
	mdStr := `# Hello, world

[aaaa](http://aaa.bbb.ccc)
`
	expected := `h1. Hello, world
[aaaa|http://aaa.bbb.ccc]`

	assert.Equal(t, expected, strings.Trim(Markdown2Confluence(mdStr), "\n"))
}

func TestSlackRequestBody(t *testing.T) {
	fmt.Println(SlackRequestBody("event", "jack", "ghost", "hello, world"))
}
