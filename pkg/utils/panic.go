package utils

import (
	"context"
	"github.com/arraisi/demogo/pkg/logger"
	"strings"

	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcPanicHandler() func(ctx context.Context, p interface{}) (err error) {
	return func(ctx context.Context, p interface{}) (err error) {
		logger.Log.WithFields(ctx,
			map[string]interface{}{
				"stack": strings.Split(strings.ReplaceAll(string(debug.Stack()), "\t", ""), "\n"),
			},
		).Errorf("RequestID: %v, error: %v", GetRequestIDFromContext(ctx), p)

		return status.Errorf(codes.Internal, "%v", "internal server error, please contact developer")
	}
}
