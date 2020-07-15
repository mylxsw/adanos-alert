package array

import "strings"

// StringUnique remove duplicate elements from array
func StringUnique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// StringsContain 判断元素是否在字符串数组中
func StringsContain(val string, items []string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}

	return false
}

// StringContainPrefix 判断字符串是否以指定的前缀开始
func StringsContainPrefix(val string, prefixs []string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(val, prefix) {
			return true
		}
	}

	return false
}

// StringsFilter 字符串数组过滤
func StringsFilter(items []string, filter func(item string) bool) []string {
	res := make([]string, 0)
	for _, item := range items {
		item = strings.TrimSpace(item)
		if filter(item) {
			res = append(res, item)
		}
	}

	return res
}

// StringsRemoveEmpty 过滤掉字符串数组中的空元素
func StringsRemoveEmpty(items []string) []string {
	return StringsFilter(items, func(item string) bool {
		return item != ""
	})
}
