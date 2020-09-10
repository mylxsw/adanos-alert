package strarr_test

import (
	"testing"

	"github.com/mylxsw/adanos-alert/pkg/strarr"
	"github.com/stretchr/testify/assert"
)

func TestStringUnique(t *testing.T) {
	arr := []string{
		"aaa",
		"bbb",
		"ccc",
		"aaa",
		"ddd",
		"ccc",
	}

	assert.EqualValues(t, 4, len(strarr.Distinct(arr)))
}

func TestStringsContainPrefix(t *testing.T) {
	s1 := "Hello, world"
	assert.True(t, strarr.HasPrefixes(s1, []string{"xxxx", "yyyy", "Hell"}))
	assert.False(t, strarr.HasPrefixes(s1, []string{"xxxx", "yyyy", "oops"}))
}

func TestStringDiff(t *testing.T) {
	itemsA := []string{"aaa", "bbb", "ccc", "ddd"}
	itemsB := []string{"ccc", "bbb", "eee"}

	res := strarr.Diff(itemsA, itemsB)
	assert.Equal(t, 2, len(res))
	assert.True(t, strarr.In("aaa", res))
	assert.True(t, strarr.In("ddd", res))
}

func TestUnion(t *testing.T) {
	itemsA := []string{"aaa", "bbb", "ccc", "ddd"}
	itemsB := []string{"ccc", "bbb", "eee"}

	res := strarr.Union(itemsA, itemsB)
	assert.Equal(t, 5, len(res))
}
