package logic

import (
	"context"

	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCaptchaLogic {
	return &GenerateCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateCaptchaLogic) GenerateCaptcha(*user.Empty) (*user.CaptchaResponse, error) {
	// todo: add your logic here and delete this line
	id, b64s, answer, err := l.svcCtx.Captcha.Generate()
	if err != nil {
		l.Logger.Errorf("[GenerateCaptchaLogic.GenerateCaptcha] error: %v", err)
		return nil, err
	}
	// save answer to cache
	if err = l.svcCtx.CaptchaStore.Set(id, answer); err != nil {
		l.Logger.Errorf("[GenerateCaptchaLogic.GenerateCaptcha] error: %v", err)
		return nil, err
	}

	return &user.CaptchaResponse{CaptchaId: id, ImgBase64: b64s}, nil
}
