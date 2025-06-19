package structutil

import (
	"errors"
	"github.com/obnahsgnaw/goutils/cacheutil"
	"github.com/obnahsgnaw/goutils/codecs/jsonutil"
	"github.com/obnahsgnaw/goutils/strutil"
	"strings"
	"time"
)

// SyncStruct 带缓存的结构体，并提供一个唯一短链（可以用自动生成的页可以手动设置一个）
type SyncStruct struct {
	cache     cacheutil.Cache
	impl      interface{}
	ttl       time.Duration
	kind      string
	shortLink string
}

func (s *SyncStruct) Init(cache cacheutil.Cache, kind string, ttl time.Duration, impl interface{}) {
	s.cache = cache
	s.kind = kind
	s.ttl = ttl
	s.impl = impl
}

func (s *SyncStruct) Save(shortLink ...string) (err error) {
	if len(shortLink) > 0 {
		s.shortLink = shortLink[0]
	}
	if s.shortLink == "" {
		s.shortLink = s.genShortLink(s.kind)
	}
	if err = s.validate(); err != nil {
		return
	}
	var data string
	if data, err = jsonutil.Encode(s.impl); err != nil {
		return
	}
	err = s.cache.Cache(s.key(), data, 0)
	return
}

func (s *SyncStruct) Load(shortLink string) (hit bool, err error) {
	s.shortLink = shortLink
	if err = s.validate(); err != nil {
		return
	}
	var data string
	if data, hit, err = s.cache.Cached(s.key()); err != nil {
		return
	}
	if !hit {
		return
	}
	err = jsonutil.Decode([]byte(data), &s.impl)
	return
}

func (s *SyncStruct) validate() (err error) {
	if s.cache == nil {
		err = errors.New("cache is nil")
		return
	}
	if s.impl == nil {
		err = errors.New("impl is nil")
		return
	}
	if s.kind == "" {
		err = errors.New("kind is empty")
		return
	}
	if s.shortLink == "" {
		err = errors.New("shortLink is empty")
		return
	}
	return
}

func (s *SyncStruct) ShortLink() string {
	return s.shortLink
}

func (s *SyncStruct) genShortLink(prefixes ...string) string {
	if s.shortLink != "" {
		return s.ShortLink()
	}
	s.shortLink = strutil.PrefixedUnique(strings.Join(prefixes, "_"))
	return s.ShortLink()
}

func (s *SyncStruct) key() string {
	return strutil.ToString(s.kind, ":", s.shortLink)
}
