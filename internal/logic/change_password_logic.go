package logic

import (
	"context"

	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"

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

// 更改密码必须通过邮箱验证才能更改
func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePassWdRequest) (*user.Empty, error) {
	// todo: change password via email verification
	if in.Code == "" || in.Email == "" {
		return nil, model.ErrVerifyEmail
	}

	if ok := l.svcCtx.MailStore.Verify(in.Email, in.Code, true); !ok {
		return nil, model.ErrVerifyEmail
	}

	if in.New == "" {
		return nil, model.ErrPasswordEmpty
	}
	// 获取用户信息
	res, err := l.svcCtx.UserModel.GetByUsernameOrEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if res == nil || res.ID == 0 {
		return nil, model.ErrUserNotFound
	}

	// 生成新密码哈希
	newHashedPwd, err := utils.EncryptPassword(in.New)
	if err != nil {
		return nil, err
	}
	// 更新用户密码
	res.Password = newHashedPwd
	res.TokenVersion = res.TokenVersion + 1

	if l.svcCtx.UserModel.Update(res) == nil {
		return nil, model.ErrChangePasswordFailed
	}

	return nil, nil
}
