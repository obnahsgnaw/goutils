package ginutil

import (
	"bytes"
	"crypto"
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/goutils/security/hsutil"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	cacheControlHeader = "Cache-Control"
	cacheControlValue  = "private, max-age=" // 缓存头
	eTagHeader         = "ETag"
	ifNoneMatchHeader  = "If-None-Match"
)

type StaticFsCache struct {
	engin         *gin.Engine
	fs            *embed.FS
	fsFallbackDir string
	rootDir       string
	etags         map[string]string
	cacheTtl      int64
	relativePath  string
	replace       map[string]func([]byte) []byte
}

func NewStaticFsCache(engin *gin.Engine, rootDir string, o ...FsCacheOption) *StaticFsCache {
	s := &StaticFsCache{
		engin:        engin,
		relativePath: "/",
		rootDir:      rootDir,
		etags:        make(map[string]string),
		replace:      make(map[string]func([]byte) []byte),
	}
	s.with(o...)
	return s
}

func (s *StaticFsCache) with(o ...FsCacheOption) {
	for _, o := range o {
		if o != nil {
			o(s)
		}
	}
}

func (s *StaticFsCache) Init() (err error) {
	var items []fs.DirEntry
	if s.fs == nil {
		items, err = os.ReadDir(s.rootDir)
	} else {
		items, err = s.fs.ReadDir(s.rootDir)
	}
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			err = s.initDir(s.rootDir, item, s.fs != nil)
		} else {
			err = s.initFile(s.rootDir, item, s.fs != nil)
		}
		if err != nil {
			return err
		}
	}

	if s.fs != nil && s.fsFallbackDir != "" {
		if err = s.initFallback(); err != nil {
			return err
		}
	}

	g := s.engin.Group(s.relativePath, gzip.Gzip(gzip.DefaultCompression), cacheMiddleware(func() *StaticFsCache { return s }, s.cacheTtl, s.relativePath))
	if s.fs != nil {
		sub, _ := fs.Sub(s.fs, s.rootDir)
		if s.fsFallbackDir == "" {
			g.StaticFS("/", http.FS(sub))
		} else {
			g.StaticFS("/", FS(sub, s.fsFallbackDir))
		}
	} else {
		g.Static("/", s.rootDir)
	}

	return nil
}

func (s *StaticFsCache) initFallback() (err error) {
	var items []fs.DirEntry
	items, err = os.ReadDir(s.fsFallbackDir)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			err = s.initDir(s.fsFallbackDir, item, false)
		} else {
			err = s.initFile(s.fsFallbackDir, item, false)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StaticFsCache) initDir(base string, entry fs.DirEntry, fsSystem bool) (err error) {
	var items []fs.DirEntry
	name := entry.Name()
	if base != "" {
		name = path.Join(base, name)
	}
	if !fsSystem {
		items, err = os.ReadDir(name)
	} else {
		items, err = s.fs.ReadDir(name)
	}
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() {
			err = s.initDir(name, item, fsSystem)
		} else {
			err = s.initFile(name, item, fsSystem)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StaticFsCache) initFile(base string, entry fs.DirEntry, fsSystem bool) (err error) {
	name := entry.Name()
	if strings.HasPrefix(name, ".") {
		return nil
	}
	if base != "" {
		name = path.Join(base, name)
	}
	var f fs.File
	if !fsSystem {
		f, err = os.Open(name)
	} else {
		f, err = s.fs.Open(name)
	}
	if err != nil {
		return err
	}
	content, err1 := io.ReadAll(f)
	if err1 != nil {
		return err1
	}
	hash, err2 := hsutil.Hash(content, crypto.SHA1)
	if err2 != nil {
		return err2
	}
	if _, ok := s.etags[name]; !ok {
		s.etags[name] = string(hash)
	}
	return nil
}

func (s *StaticFsCache) etag(filename string) string {
	if etag, ok := s.etags[filename]; ok {
		return etag
	}
	return ""
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w CustomResponseWriter) RawWrite(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func cacheMiddleware(s func() *StaticFsCache, ttl int64, prefix string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ttl > 0 {
			if prefix == "" || strings.HasPrefix(c.Request.URL.Path, prefix) {
				// 设置缓存控制头
				c.Header(cacheControlHeader, cacheControlValue+strconv.FormatInt(ttl, 10))

				// 生成并设置 ETag 头
				rqPath := c.Request.URL.Path
				if s().fs != nil {
					rqPath = path.Join(s().rootDir, rqPath)
				}
				eTag := s().etag(rqPath)
				c.Header(eTagHeader, eTag)

				// 检查 If-None-Match 头与生成的 ETag 是否匹配，若匹配则返回 304 Not Modified
				if match := c.GetHeader(ifNoneMatchHeader); match != "" {
					if match == eTag {
						c.Status(http.StatusNotModified)
						c.Abort()
						return
					}
				}
			}
		}
		crw := &CustomResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = crw
		c.Next()
		if rp, ok := s().replace[c.Request.URL.Path]; ok {
			_, _ = crw.RawWrite(rp(crw.body.Bytes()))
		} else {
			_, _ = crw.RawWrite(crw.body.Bytes())
		}
	}
}
