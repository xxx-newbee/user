package mail

import (
	"errors"
	"fmt"

	"github.com/xxx-newbee/storage"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/logic/utils"
	"gopkg.in/gomail.v2"
)

const (
	prefix     = "user:mail"
	RedisTopic = "mail"
)

type (
	MailVerifyStore struct {
		dialer *gomail.Dialer
		cache  storage.AdapterCache
		expire int
	}

	MailStore interface {
		Set(mail, value string) error
		Get(mail string, clear bool) string
		Verify(mail, value string, clear bool) bool
		Consumer(msg storage.Messager) error
	}
)

func NewMailVerifyStore(c config.Config, cache storage.AdapterCache, expire int) MailStore {
	return &MailVerifyStore{
		dialer: gomail.NewDialer(c.SMTP.Host, c.SMTP.Port, c.SMTP.MailFrom, c.SMTP.Password),
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

func (m *MailVerifyStore) Consumer(messager storage.Messager) error {
	// 生成验证码
	code, err := utils.GenerateCode(6)
	if err != nil {

		//lock.Release(context.TODO())
		return err
	}
	to, ok := messager.GetValues()["email"].(string)
	if !ok {
		return errors.New("email convert string error")
	}
	// 生成邮件文本
	message := gomail.NewMessage()
	message.SetHeader("From", m.dialer.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "邮箱验证 - 验证码")
	body := fmt.Sprintf(`
	<html>
		<body>
			<h2>邮箱验证</h2>
			<p>您好，</p>
			<p>感谢您使用我们的服务。您的验证码是：<strong>%s</strong></p>
			<p>此验证码将在5分钟内过期，请尽快使用。</p>
			<br>
			<p>如果您没有请求此验证码，请忽略此邮件。</p>
		</body>
	</html>
	`, code)
	message.SetBody("text/html", body)

	// 发送邮件
	err = m.dialer.DialAndSend(message)

	if err != nil {
		//_ = lock.Release(context.TODO())
		return err
	}

	err = m.Set(to, code)
	if err != nil {
		//lock.Release(context.TODO())
		return err
	}

	return nil
}
