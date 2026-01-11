package logic

import (
	"context"
	"errors"

	"github.com/xxx-newbee/user/internal/dao"
	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReqest) (*user.UpdateUserInfoResponse, error) {
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

	if in.Nickname != "" {
		res.Nickname = in.Nickname
	}

	// 要更改钱包地址和密码还需要进行密码进一步去验证，仅依靠token验证不安全
	if in.WalletAddr != "" {
		res.Wallet = in.WalletAddr
	}
	if in.Password != "" {
		hashedPassword, err := utils.EncryptPassword(in.Password)
		if err != nil {
			return nil, err
		}
		res.Password = hashedPassword
	}

	userDao.Update(res)
	return &user.UpdateUserInfoResponse{}, nil
}
