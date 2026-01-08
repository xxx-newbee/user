package logic

import (
	"context"

	"github.com/xxx-newbee/user/internal/dao"
	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"

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
	user_dao := dao.NewUserDao()
	res, err := user_dao.GetByUsername(in.Username)
	if err != nil {
		return nil, err
	}
	if res == nil || res.ID == 0 {
		return nil, model.ErrUsernameOrPasswordIncorrect
	}

	if in.Password == "" {
		return nil, model.ErrPasswordNeccessary
	}

	if err := utils.ComparePassword(res.Password, in.Password); err != nil {
		return nil, model.ErrUsernameOrPasswordIncorrect
	}

	// jwt token generation can be added here
	token, err := utils.GenerateJWTToken(int64(res.ID), res.Username, l.svcCtx.Config.JWT.Secret)
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
