package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *GetUserInfoLogic) GetUserInfo(*user.Empty) (*user.GetUserInfoResponse, error) {
	MD, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errors.New("metadata not found in context")
	}

	tokenStrs := MD.Get("authorization")

	if len(tokenStrs) == 0 || (len(tokenStrs) == 1 && tokenStrs[0] == "") {
		return nil, errors.New("illegal usage")
	}
	tokenStr := tokenStrs[0]
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	claims, err := utils.ParseJWTToken(tokenStr, l.svcCtx.Config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	username := claims.Username

	res, err := model.GetByUsername(l.svcCtx.Database, username)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, model.ErrUserNotFound
	}
	if res.TokenVersion != claims.TokenVersion {
		return nil, model.ErrTokenExpired
	}

	return &user.GetUserInfoResponse{
		Id:               int64(res.ID),
		Username:         res.Username,
		Nickname:         res.Nickname,
		WalletAddr:       res.Wallet,
		UserReferralCode: res.UserReferralCode,
		ReferralCode:     res.ReferralCode,
	}, nil

	return &user.GetUserInfoResponse{}, nil
}
