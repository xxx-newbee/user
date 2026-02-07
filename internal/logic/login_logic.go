package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/xxx-newbee/storage"
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
	var username = in.Username
	var status = "2"
	var msg = "登录成功"
	// 日志入库
	defer l.LoginLogToDB(username, status, msg)
	// 检查验证码
	ck := l.svcCtx.CaptchaStore.Verify(in.CaptchaId, in.CaptchaCode, true)
	if ck != true {
		status = "1"
		msg = model.ErrCaptchaIncorrect.Error()
		return nil, model.ErrCaptchaIncorrect
	}
	// 获取用户
	res, err := model.GetByUsername(l.svcCtx.Database, in.Username)
	if err != nil {
		status = "1"
		msg = err.Error()
		return nil, err
	}
	if res == nil || res.ID == 0 {
		status = "1"
		msg = "username not found"
		return nil, model.ErrUsernameOrPasswordIncorrect
	}

	if in.Password == "" {
		status = "1"
		msg = "password empty"
		return nil, model.ErrPasswordNecessary
	}

	if err := utils.ComparePassword(res.Password, in.Password); err != nil {
		status = "1"
		msg = "password incorrect"
		return nil, model.ErrUsernameOrPasswordIncorrect
	}

	// jwt token generation can be added here
	token, err := utils.GenerateJWTToken(int64(res.ID), res.Username, l.svcCtx.Config.JWT.Secret, res.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		UserId:           int64(res.ID),
		Token:            token,
		Username:         res.Username,
		Nickname:         res.Nickname,
		WalletAddr:       res.Wallet,
		UserReferralCode: res.UserReferralCode,
		ReferralCode:     res.ReferralCode,
	}, nil
}

// 登录日志消费者
func (l *LoginLogic) SaveLoginLog(msg storage.Messager) error {
	rb, err := json.Marshal(msg.GetValues())
	if err != nil {
		l.Logger.Error(err.Error())
		return err
	}
	var ll model.SysLoginLog
	if err = json.Unmarshal(rb, &ll); err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	if err = l.svcCtx.Database.Create(&ll).Error; err != nil {
		l.Logger.Error(err.Error())
		return err
	}

	return nil
}

// 登录日志生产者
func (l *LoginLogic) LoginLogToDB(username, status, msg string) {
	ll := make(map[string]interface{})
	MD, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		l.Logger.Error("metadata not found in context")
	}

	if uas := MD.Get("UA"); len(uas) > 0 {
		logx.Infof("user-agent: %s", uas)
		ua := user_agent.New(uas[0])
		ll["remark"] = uas[0]
		ll["os"] = ua.OS()
		browserName, browserVersion := ua.Browser()
		ll["browser"] = browserName + " " + browserVersion
		ll["platform"] = ua.Platform()
	}
	if ip_addr := MD.Get("remote-addr"); len(ip_addr) > 0 {
		logx.Infof("remote-addr: %s", ip_addr)
		ll["ipaddr"] = ip_addr[0]
	}

	ll["loginTime"] = time.Now()
	ll["status"] = status
	ll["msg"] = msg
	ll["username"] = username

	// 创建消息
	message := &queue.Message{
		Stream: model.SysLoginLog{}.TableName(),
		Values: ll,
	}
	// 消息入队
	if err := l.svcCtx.MemoryQueue.Append(message); err != nil {
		l.Logger.Errorf("Append message error: %s", err.Error())
	}

}
