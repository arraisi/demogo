package server

import (
	"context"
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/pkg/constant"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"time"
)

func UnaryRequestInterceptor(conf *config.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := uuid.NewString()
		reqIdCtx := context.WithValue(ctx, constant.REQUEST_ID, requestID)

		_, ok := ctx.Deadline()
		if ok {
			return handler(reqIdCtx, req)
		}

		ctxDefault, cancel := context.WithTimeout(reqIdCtx, (conf.DefaultTimeout * time.Second))
		defer cancel()

		return handler(ctxDefault, req)
	}
}
