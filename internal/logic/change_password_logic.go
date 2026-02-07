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

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePassWdRequest) (*user.Empty, error) {
	if in.Old == "" || in.New == "" {
		return nil, model.ErrPasswordEmpty
	}
	// 元数据获取token，解析token获取用户名
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
	// 获取用户信息
	//res, err := dao.NewUserDao().GetByUsername(username)
	res, err := model.GetByUsername(l.svcCtx.Database, username)
	if err != nil {
		return nil, err
	}
	if res == nil || res.ID == 0 {
		return nil, model.ErrUserNotFound
	}
	// 旧密码匹配
	if err := utils.ComparePassword(res.Password, in.Old); err != nil {
		return nil, model.ErrOldPasswordIncorrect
	}
	// 生成新密码哈希
	newHashedPwd, err := utils.EncryptPassword(in.New)
	if err != nil {
		return nil, err
	}
	// 更新用户密码
	res.Password = newHashedPwd
	res.TokenVersion = res.TokenVersion + 1

	if model.UpdateUser(l.svcCtx.Database, res) == nil {
		return nil, model.ErrChangePasswordFailed
	}

	return nil, nil
}
