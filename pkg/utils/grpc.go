package utils

import (
	"context"
	"crypto/tls"
	"demogo/config"
	"fmt"
	"google.golang.org/grpc/codes"
	"time"

	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateGRPCConnection(conf config.GRPCItemConfig) (conn *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Timeout)*time.Second)
	defer cancel()

	var opts []grpc.DialOption
	cred := insecure.NewCredentials()
	if conf.TLS {
		configTls := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS13,
		}
		cred = credentials.NewTLS(configTls)
	}
	opts = append(opts, grpc.WithTransportCredentials(cred))
	if conf.EnableRetry {
		exponentialBackoff := time.Duration(conf.ExponentialBackoffRetry) * time.Millisecond
		retryOpts := grpc.WithUnaryInterceptor(
			grpcRetry.UnaryClientInterceptor(
				grpcRetry.WithCodes(codes.Aborted, codes.Unavailable, codes.Unknown),
				grpcRetry.WithMax(uint(conf.MaxRetry)),
				grpcRetry.WithBackoff(grpcRetry.BackoffExponential(exponentialBackoff)),
			),
		)
		opts = append(opts, retryOpts)
	}

	connUrl := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	return grpc.DialContext(ctx, connUrl, opts...)
}
