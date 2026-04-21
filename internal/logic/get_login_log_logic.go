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

type GetLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogLogic {
	return &GetLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoginLogLogic) GetLoginLog(in *user.GetLoginLogRequest) (*user.LoginLogResponse, error) {

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

	logs, err := l.svcCtx.SysLoginLog.GetUserLoginLogs(user_id, int(in.Page))
	if err != nil {
		return nil, err
	}
	resp := &user.LoginLogResponse{}
	for _, log := range logs {
		resp.LoginLogs = append(resp.LoginLogs, &user.LoginLog{
			UserId:    log.UserId,
			Status:    log.Status,
			IpAddr:    log.Ipaddr,
			Location:  log.LoginLocation,
			Browser:   log.Browser,
			Os:        log.Os,
			Platform:  log.Platform,
			LoginTime: log.LoginTime.Unix(),
			Remark:    log.Remark,
			Msg:       log.Msg,
		})
	}

	return resp, nil
}
