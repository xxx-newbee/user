package logic

import (
	"context"
	"errors"

	"github.com/xxx-newbee/go-micro/user/internal/dao"
	"github.com/xxx-newbee/go-micro/user/internal/logic/utils"
	"github.com/xxx-newbee/go-micro/user/internal/svc"
	"github.com/xxx-newbee/go-micro/user/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errors.New("metadata not found in context")
	}

	tokenStrs := md.Get("Authorization")
	if len(tokenStrs) == 0 {
		return nil, errors.New("illegal usage")
	}

	tokenStr := tokenStrs[0]
	claims, err := utils.ParseJWTToken(tokenStr, l.svcCtx.Config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	userId := claims.UserID

	userDao := dao.NewUserDao()
	res, err := userDao.GetByID(uint(userId))
	if err != nil {
		return nil, err
	}

	return &user.GetUserInfoResponse{
		Id:               int64(res.ID),
		Username:         res.Username,
		Nickname:         res.Nickname,
		WalletAddr:       res.Wallet,
		UserReferralCode: res.UserReferralCode,
		ReferralCode:     res.ReferralCode,
	}, nil
}
