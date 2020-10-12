package misc

import (
	"strings"

	"github.com/pingcap/parser"
)

// SQLFinger 生成 SQL 指纹
func SQLFinger(sqlStr string) string {
	return strings.ReplaceAll(parser.Normalize(sqlStr), " . ", ".")
}
