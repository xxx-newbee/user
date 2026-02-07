package captcha

import (
	"fmt"

	"github.com/mojocn/base64Captcha"
	"github.com/xxx-newbee/storage"
)

const (
	prefix = "user:captcha"
)

type CaptchaStore struct {
	cache  storage.AdapterCache
	expire int
}

func NewCaptchaStore(cache storage.AdapterCache, expire int) base64Captcha.Store {
	return &CaptchaStore{
		cache:  cache,
		expire: expire,
	}
}

func (c *CaptchaStore) Set(id, value string) error {
	key := fmt.Sprintf("%s:%s", prefix, id)
	return c.cache.Set(key, value, c.expire)
}

func (c *CaptchaStore) Get(id string, clear bool) string {
	key := fmt.Sprintf("%s:%s", prefix, id)
	value, err := c.cache.Get(key)
	if err != nil {
		return ""
	}
	if clear {
		_ = c.cache.Del(key)
	}
	return value
}

func (c *CaptchaStore) Verify(id, ans string, clear bool) bool {
	if id == "" || ans == "" {
		return false
	}
	return c.Get(id, clear) == ans

}
