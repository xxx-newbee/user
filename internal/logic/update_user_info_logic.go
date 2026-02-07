package logic

import (
	"context"
	"errors"

	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
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

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReqest) (*user.Empty, error) {
	// todo: add your logic here and delete this line
	MD, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errors.New("metadata not found in context")
	}

	tokenStrs := MD.Get("Authorization")
	if len(tokenStrs) == 0 {
		return nil, errors.New("illegal usage")
	}

	tokenStr := tokenStrs[0]
	claims, err := utils.ParseJWTToken(tokenStr, l.svcCtx.Config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	username := claims.Username

	//userDao := dao.NewUserDao()
	//res, err := userDao.GetByUsername(username)
	res, err := model.GetByUsername(l.svcCtx.Database, username)
	if err != nil {
		return nil, err
	}
	if res == nil || res.ID == 0 {
		return nil, model.ErrUserNotFound
	}
	if res.TokenVersion != claims.TokenVersion {
		return nil, model.ErrTokenExpired
	}
	if in.Nickname != "" {
		res.Nickname = in.Nickname
	}

	// 要更改钱包地址还需要进行密码进一步去验证，仅依靠token不安全
	if in.WalletAddr != "" {
		res.Wallet = in.WalletAddr
	}

	//if err = userDao.Update(res); err != nil {
	//	return nil, model.ErrUpdateUserFailed
	//}
	if model.UpdateUser(l.svcCtx.Database, res) == nil {
		return nil, model.ErrUpdateUserFailed
	}
	return nil, nil
}
