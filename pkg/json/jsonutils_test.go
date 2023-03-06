package json_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/json"
)

var message1 = `{
    "message": "ack_confirm",
    "context": {
        "msg": "层级消息内容",
        "sms": {
            "id": 44444,
            "app_id": 1,
            "template_params": {
                "username": "李逍遥",
                "password": "lixiaoyao",
                "gender": null,
                "created_at": "2018-11-12 13:47:55"
            },
            "status": 1,
            "created_at": "2018-11-12 13:47:55",
            "updated_at": "2018-11-14 13:49:04"
        },
        "ack": {
            "msg": "😄",
            "code": "6460"
        },
        "file": "/webroot/your/project/Test.php:322"
    },
    "level": null,
    "level_name": "ERROR",
    "channel": "custom_cmd",
    "datetime": "2018-11-16 13:51:01",
    "extra": {
        "ref": "5bee5ac564a71bbb33cai2jkk"
    }
}`

var message2 = `{
	"message": null,
	"context": {
		"msg": null,
		"reason": "unknown",
		"extra": {
			"numbers": [],
			"numbers2": [1, 2, 3, 4, 5],
			"users": ["user1", "user2"],
			"mix": ["string1", 45],
			"nulls": [null, null, 5],
			"user": {}
		}
	}
}`

var message3 = `{
	"messages": [
		{
			"id": 123,
			"key": "xxx"
		},
		{
			"id": 444,
			"key": "yyyy"
		}
	]
}`

var message4 = `[
	{
		"key": 123
	},
	{
		"key": 444
	}
]`

var message5 = `[13, 44, 55, 66]`

func TestToKvPairs(t *testing.T) {
	ju, err := json.New([]byte(message1), 0, false)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	kvPairs := ju.ToKvPairs()
	if len(kvPairs) == 0 {
		t.Error("convert to kv pairs failed")
	}

	if len(kvPairs) != 19 {
		t.Errorf("kv pairs not matched")
	}

}

func TestToKvPairsArray(t *testing.T) {
	ju, err := json.New([]byte(message1), 0, false)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	kvPairs := ju.ToKvPairsArray()
	if len(kvPairs) != 19 {
		t.Errorf("kv pairs not matched")
	}
}

func TestNullValue(t *testing.T) {
	ju, err := json.New([]byte(message2), 0, false)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairs()
	if v, ok := pairs["context.msg"]; !ok || v != "(null)" {
		t.Errorf("kv pairs with null value test failed")
	}

	if v, ok := pairs["message"]; !ok || v != "(null)" {
		t.Errorf("kv pairs with null value test failed")
	}

	// for k, v := range pairs {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
}

func TestKvPairsWithLevelLimit(t *testing.T) {
	ju, err := json.New([]byte(message1), 2, false)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairsArray()
	for _, kv := range pairs {
		// fmt.Printf("%s : %s\n", kv.Key, kv.Value)
		if len(strings.Split(kv.Key, ".")) > 2 {
			t.Error("test kv pairs with level limit failed")
		}
	}
}
func TestNullValueSkipSimpleValue(t *testing.T) {
	ju, err := json.New([]byte(message2), 0, true)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairs()
	if v, ok := pairs["context.msg"]; !ok || v != "(null)" {
		t.Errorf("kv pairs with null value test failed")
	}

	if v, ok := pairs["message"]; !ok || v != "(null)" {
		t.Errorf("kv pairs with null value test failed")
	}

	// for k, v := range pairs {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
}

func TestComplexArrayValue(t *testing.T) {
	ju, err := json.New([]byte(message3), 2, true)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairsArray()
	// for _, kv := range pairs {
	// 	fmt.Printf("%s : %s\n", kv.Key, kv.Value)
	// }

	if len(pairs) != 2 {
		t.Errorf("kv pairs test failed for complex array")
	}
}

func TestRootArrayType(t *testing.T) {
	ju, err := json.New([]byte(message4), 2, true)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairsArray()

	for _, kv := range pairs {
		fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	}
}

func TestRootArrayType2(t *testing.T) {
	ju, err := json.New([]byte(message5), 1, true)
	if err != nil {
		t.Errorf("parse json failed: %s", err.Error())
	}

	pairs := ju.ToKvPairsArray()

	// for _, kv := range pairs {
	// 	fmt.Printf("%s: %s\n", kv.Key, kv.Value)
	// }

	if len(pairs) <= 0 {
		t.Errorf("test failed")
	}
}
