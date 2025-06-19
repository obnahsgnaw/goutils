package strutil

import (
	"github.com/obnahsgnaw/application/pkg/utils"
	"github.com/rs/xid"
)

func Unique() string {
	return xid.New().String()
}

func PrefixedUnique(prefix string) string {
	if prefix == "" {
		return Unique()
	}
	return utils.ToStr(prefix, Unique())
}
