package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/lockerutil"
	"strconv"
	"time"
)

const dLockerKeyPrefix = "distributed-locker"

type Builder struct {
	client *redis.Client
}

func New(client *redis.Client) *Builder {
	return &Builder{client: client}
}

func (b *Builder) Compete(key string, timeout time.Duration) lockerutil.Locker {
	return newLocker(b.client, key, timeout)
}

type Locker struct {
	client    *redis.Client
	err       error
	exist     bool
	key       string
	val       string
	createdAt time.Time
	ttl       time.Duration
}

func lockerErr(msg string) error {
	return errutil.New("redis locker util: ", msg)
}

func newLocker(client *redis.Client, key string, ttl time.Duration) *Locker {
	now := time.Now()
	val := strconv.FormatInt(now.UnixNano(), 10)
	key = dLockerKeyPrefix + ":" + key
	s := &Locker{
		key:       key,
		val:       val,
		createdAt: now,
		ttl:       ttl,
		client:    client,
	}
	if rs := client.SetNX(context.Background(), key, val, ttl); rs.Err() != nil {
		s.err = lockerErr(rs.Err().Error())
	} else {
		if !rs.Val() {
			s.exist = true
		}
	}
	return s
}

func (l *Locker) Unlock() {
	if l.Error() == nil {
		if l.createdAt.Add(l.ttl).Before(time.Now()) {
			return
		}
		sc := `if redis.call("get",KEYS[1]) == ARGV[1] then return redis.call("del",KEYS[1]) else return 0 end`
		_ = l.client.Eval(context.Background(), sc, []string{l.key}, l.val)
	}
}

func (l *Locker) Error() error {
	return l.err
}

func (l *Locker) Hit() bool {
	return !l.exist
}
