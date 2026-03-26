package logic

import (
	"context"
	"fmt"

	"github.com/xxx-newbee/storage/queue"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/internal/svc/mail"
	"github.com/xxx-newbee/user/user"
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
	_, err := l.svcCtx.Locker.Lock(lockKey, 60, nil) // 2分钟邮件锁
	if err != nil {
		// 未到冷却时间
		l.Logger.Errorf(" 📧 email locker err: %s", err.Error())
		return &user.SendEmailResponse{Success: false}, nil
	}

	l.ToQueue(in.Email)
	return &user.SendEmailResponse{Success: true}, nil
}

func (l *SendEmailLogic) ToQueue(email string) {
	message := &queue.Message{
		Stream: mail.RedisTopic,
		Values: map[string]interface{}{
			"email": email,
		},
	}
	if err := l.svcCtx.RedisQueue.Append(message); err != nil {
		l.Logger.Errorf("Append email message err: %s", err.Error())
	}
}
