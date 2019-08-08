package matcher

import (
	"strings"
)

type Helpers struct{}

// Lower returns a copy of the string s with all Unicode letters mapped to their lower case.
func (Helpers) Lower(val string) string {
	return strings.ToLower(val)
}

// Upper returns a copy of the string s with all Unicode letters mapped to their upper case.
func (Helpers) Upper(val string) string {
	return strings.ToUpper(val)
}
