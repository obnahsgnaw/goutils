package strutil

import (
	"github.com/rs/xid"
)

func Unique() string {
	return xid.New().String()
}

func PrefixedUnique(prefix string) string {
	if prefix == "" {
		return Unique()
	}
	return ToString(prefix, Unique())
}
