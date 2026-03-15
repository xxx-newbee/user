package mail

import (
	"fmt"

	"github.com/xxx-newbee/storage"
)

const (
	prefix = "user:mail"
)

type (
	MailVerifyStore struct {
		cache  storage.AdapterCache
		expire int
	}

	MailStore interface {
		Set(mail, value string) error
		Get(mail string, clear bool) string
		Verify(mail, value string, clear bool) bool
	}
)

func NewMailVerifyStore(cache storage.AdapterCache, expire int) MailStore {
	return &MailVerifyStore{
		cache:  cache,
		expire: expire,
	}
}

func (m *MailVerifyStore) Set(mail, value string) error {
	key := fmt.Sprintf("%s:%s", prefix, mail)
	return m.cache.Set(key, value, m.expire)
}

func (m *MailVerifyStore) Get(mail string, clear bool) string {
	key := fmt.Sprintf("%s:%s", prefix, mail)
	value, err := m.cache.Get(key)
	if err != nil {
		return ""
	}
	if clear {
		_ = m.cache.Del(key)
	}
	return value
}

func (m *MailVerifyStore) Verify(mail, value string, clear bool) bool {
	if mail == "" || value == "" {
		return false
	}
	return m.Get(mail, clear) == value
}
