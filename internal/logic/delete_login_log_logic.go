package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/xxx-newbee/user/internal/logic/utils"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoginLogLogic {
	return &DeleteLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLoginLogLogic) DeleteLoginLog(in *user.DeleteLoginLogRequest) (*user.Empty, error) {

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

	user_id := claims.UserID
	log, err := l.svcCtx.SysLoginLog.GetLoginLogById(in.Id)
	if err != nil {
		return nil, err
	}
	// 校验用户名
	if user_id == log.UserId {
		err := l.svcCtx.SysLoginLog.DeleteLoginLog(in.Id)
		return &user.Empty{}, err
	}
	return &user.Empty{}, nil
}
