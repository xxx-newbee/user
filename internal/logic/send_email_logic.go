package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/bsm/redislock"
	"github.com/xxx-newbee/storage"
	"github.com/xxx-newbee/storage/queue"
	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"gopkg.in/gomail.v2"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailLogic {
	return &SendEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendEmailLogic) SendEmail(in *user.SendEmailRequest) (*user.SendEmailResponse, error) {
	// 先检查分布式锁
	lockKey := fmt.Sprintf("email:locker:%s", in.Email)
	lock, err := l.svcCtx.Locker.Lock(lockKey, 60, nil) // 2分钟邮件锁
	if err != nil {
		// 未到冷却时间
		l.Logger.Errorf(" 📧 email locker err: %s", err.Error())
		return &user.SendEmailResponse{Success: false}, nil
	}

	l.ToQueue(in.Email, lock)
	return &user.SendEmailResponse{Success: true}, nil
}

func (l *SendEmailLogic) ToQueue(email string, locker *redislock.Lock) {
	message := &queue.Message{
		Stream: "email",
		Values: map[string]interface{}{
			"email":  email,
			"locker": locker,
		},
	}
	if err := l.svcCtx.MemoryQueue.Append(message); err != nil {
		l.Logger.Errorf("Append email message err: %s", err.Error())
	}
}

func (l *SendEmailLogic) MailComsumer(messager storage.Messager) error {
	lock, ok := messager.GetValues()["locker"].(*redislock.Lock)
	if !ok {
		return errors.New("locker is not redislock")
	}
	ttl, err := lock.TTL(context.TODO())
	if ttl < 0 || err != nil {
		return fmt.Errorf("locker ttl %f, err: %s", ttl.Seconds(), err.Error())
	}
	// 生成验证码
	code, err := utils.GenerateCode(6)
	if err != nil {
		l.Logger.Errorf(" 📧 generate code err: %s", err.Error())
		lock.Release(context.TODO())
		return err
	}
	to, ok := messager.GetValues()["email"].(string)
	if !ok {
		return errors.New("email convert string error")
	}
	// 生成邮件文本
	message := gomail.NewMessage()
	message.SetHeader("From", l.svcCtx.Config.SMTP.MailFrom)
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
	err = l.svcCtx.Mail.DialAndSend(message)

	if err != nil {
		l.Logger.Errorf(" 📧 send email err: %s", err.Error())
		_ = lock.Release(context.TODO())
		return err
	}

	err = l.svcCtx.MailStore.Set(to, code)
	if err != nil {
		l.Logger.Errorf(" 📧 store code err: %s", err.Error())
		lock.Release(context.TODO())
		return err
	}

	return nil
}
