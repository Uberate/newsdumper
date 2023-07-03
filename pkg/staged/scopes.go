package staged

import "strings"

const (
	InnerScope = "__inner"
)

var SharedMemoryScope = Scopes(InnerScope, "shared")

func Scopes(scopes ...string) string {
	return strings.Join(scopes, ":")
}
