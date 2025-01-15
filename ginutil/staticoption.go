package ginutil

import (
	"embed"
	"strings"
)

type FsCacheOption func(*StaticFsCache)

func CaCheTtl(ttl int64) FsCacheOption {
	return func(s *StaticFsCache) {
		s.cacheTtl = ttl
	}
}

func Fs(fs *embed.FS) FsCacheOption {
	return func(s *StaticFsCache) {
		s.fs = fs
	}
}

func RelativePath(relativePath string) FsCacheOption {
	return func(s *StaticFsCache) {
		s.relativePath = "/" + strings.TrimPrefix(relativePath, "/")
	}
}

func Replace(rp map[string]func([]byte) []byte) FsCacheOption {
	return func(s *StaticFsCache) {
		if rp != nil {
			s.replace = rp
		}
	}
}
