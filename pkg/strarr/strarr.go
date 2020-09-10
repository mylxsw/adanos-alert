package strarr

import "strings"

// Distinct remove duplicate elements from array
func Distinct(input []string) []string {
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

// In 判断元素是否在字符串数组中
func In(val string, items []string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}

	return false
}

// HasPrefixes 判断字符串是否以指定的前缀开始
func HasPrefixes(val string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(val, prefix) {
			return true
		}
	}

	return false
}

// Filter 字符串数组过滤
func Filter(items []string, filter func(item string) bool) []string {
	res := make([]string, 0)
	for _, item := range items {
		item = strings.TrimSpace(item)
		if filter(item) {
			res = append(res, item)
		}
	}

	return res
}

// FilterEmpty 过滤掉字符串数组中的空元素
func FilterEmpty(items []string) []string {
	return Filter(items, func(item string) bool {
		return item != ""
	})
}

// Diff 提取 itemsA 中包含，但是 itemsB 中不存在的元素
func Diff(itemsA []string, itemsB []string) []string {
	res := make([]string, 0)
	for _, item := range itemsA {
		if In(item, itemsB) {
			continue
		}

		res = append(res, item)
	}

	return res
}

// Union 两个字符串数组合并，去重复
func Union(itemsA []string, itemsB []string) []string {
	return Distinct(append(itemsA, itemsB...))
}
