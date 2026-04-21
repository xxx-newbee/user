package logic

import (
	"context"
	"strings"
	"time"

	"github.com/xxx-newbee/storage/queue"
	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"google.golang.org/grpc/metadata"

	"github.com/mssola/user_agent"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	var user_id = 0
	var status = "2"
	var msg = "登录成功"
	// 日志入库
	defer func() {
		l.LoginLogToQueue(uint64(user_id), status, msg)
	}()

	// 检查验证码
	ck := l.svcCtx.CaptchaStore.Verify(in.CaptchaId, in.CaptchaCode, true)
	if ck != true {
		status = "1"
		msg = model.ErrCaptchaIncorrect.Error()
		return nil, model.ErrCaptchaIncorrect
	}

	// 获取用户
	res, err := l.svcCtx.UserModel.GetByUsernameOrEmail(in.Username)
	if err != nil {
		status = "1"
		msg = "错误"
		return nil, err
	}
	if res == nil || res.ID == 0 {
		status = "1"
		msg = "未找到用户"
		return nil, model.ErrUsernameOrPasswordIncorrect
	}
	user_id = int(res.ID)

	if in.Password == "" {
		status = "1"
		msg = "密码错误"
		return nil, model.ErrPasswordNecessary
	}

	if err := utils.ComparePassword(res.Password, in.Password); err != nil {
		status = "1"
		msg = "密码错误"
		return nil, model.ErrUsernameOrPasswordIncorrect
	}

	// jwt token generation can be added here
	token, err := utils.GenerateJWTToken(uint64(res.ID), res.Username, l.svcCtx.Config.JWT.Secret, res.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		UserId:           int64(res.ID),
		Token:            token,
		Username:         res.Username,
		Nickname:         res.Nickname,
		Email:            res.Email,
		WalletAddr:       res.Wallet,
		UserReferralCode: res.UserReferralCode,
		ReferralCode:     res.ReferralCode,
	}, nil
}

// 登录日志生产者
func (l *LoginLogic) LoginLogToQueue(userId uint64, status, msg string) {
	ll := make(map[string]interface{})
	MD, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		l.Logger.Error("metadata not found in context")
	}

	if uas := MD.Get("UA"); len(uas) > 0 {
		ua := user_agent.New(uas[0])
		ll["remark"] = uas[0]
		ll["os"] = ua.OS()
		browserName, browserVersion := ua.Browser()
		ll["browser"] = browserName + " " + browserVersion
		ll["platform"] = ua.Platform()
	}
	if ip_addr := MD.Get("remote-addr"); len(ip_addr) > 0 {
		// ip_addr[0] = 93.179.101.214:33866
		ll["ipaddr"] = strings.Split(ip_addr[0], ":")[0]
	}

	ll["loginTime"] = time.Now()
	ll["status"] = status
	ll["msg"] = msg
	ll["userId"] = userId
	l.Logger.Infof("user_id: %d", userId)
	// 创建消息
	message := &queue.Message{
		Stream: model.SysLoginLogModel{}.TableName(),
		Values: ll,
	}
	// 消息入队
	if err := l.svcCtx.MemoryQueue.Append(message); err != nil {
		l.Logger.Errorf("Append login log message error: %s", err.Error())
	}

}
