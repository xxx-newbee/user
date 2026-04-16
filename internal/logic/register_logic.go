package logic

import (
	"context"

	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, model.ErrUsernameOrPasswordEmpty
	}

	// 邮箱验证
	if in.Email != "" {
		if ok := l.svcCtx.MailStore.Verify(in.Email, in.EmailCode, true); !ok {
			return nil, model.ErrVerifyEmail
		}
	}

	// 校验验证码
	ck := l.svcCtx.CaptchaStore.Verify(in.CaptchaId, in.CaptchaCode, true)
	if ck != true {
		return nil, model.ErrCaptchaIncorrect
	}

	res, err := l.svcCtx.UserModel.GetByUsernameOrEmail(in.Username)
	if err != nil {
		return nil, err
	}
	if res != nil && res.ID > 0 {
		return nil, model.ErrUserAlreadyExist
	}

	hashedPassword, err := utils.EncryptPassword(in.Password)
	if err != nil {
		return nil, err
	}

	referralCode, err := utils.GenerateReferralCode()
	if err != nil {
		return nil, model.ErrGenerateReferralCode
	}

	newUser := &model.User{
		Username:         in.Username,
		Password:         hashedPassword,
		Nickname:         in.Nickname,
		Email:            in.Email,
		Wallet:           in.WalletAddr,
		UserReferralCode: referralCode,
		ReferralCode:     in.ReferralCode,
		TokenVersion:     0,
	}

	if err := l.svcCtx.UserModel.Create(newUser); err != nil {
		return nil, model.ErrUserCreateFailed
	}

	return &user.RegisterResponse{
		Username:         newUser.Username,
		Nickname:         newUser.Nickname,
		Email:            newUser.Email,
		UserReferralCode: newUser.UserReferralCode,
		ReferralCode:     newUser.ReferralCode,
		WalletAddr:       newUser.Wallet,
	}, nil
}
