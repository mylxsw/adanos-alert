package json

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

// Gets 从json中提取单个值，可以使用逗号分割多个key作为备选
func Gets(key string, defaultValue string, body string) string {
	keys := strings.Split(key, ",")
	var res string
	for _, k := range keys {
		res = Get(k, "", body)
		if res != "" {
			return res
		}
	}

	return defaultValue
}

// Get 从json中提取单个值
func Get(key string, defaultValue string, body string) string {
	keys := strings.Split(key, ".")

	value, dataType, _, err := jsonparser.Get([]byte(body), keys...)
	if err != nil {
		return defaultValue
	}

	switch dataType {
	case jsonparser.NotExist:
		return defaultValue
	case jsonparser.String:
		if res, err := jsonparser.ParseString(value); err == nil {
			return res
		}
	case jsonparser.Number:
		if res, err := jsonparser.ParseFloat(value); err == nil {
			return strconv.FormatFloat(res, 'f', -1, 32)
		}
		if res, err := jsonparser.ParseInt(value); err == nil {
			return fmt.Sprintf("%d", res)
		}
	case jsonparser.Object:
		fallthrough
	case jsonparser.Array:
		return fmt.Sprintf("%s", value)
	case jsonparser.Boolean:
		if res, err := jsonparser.ParseBoolean(value); err == nil {
			if res {
				return "true"
			} else {
				return "false"
			}
		}
	case jsonparser.Null:
		return "null"
	case jsonparser.Unknown:
		return "unknown"
	}

	return defaultValue
}
