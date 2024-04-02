package server

import (
	"context"
	"demogo/config"
	productGrpc "demogo/internal/app/api/grpc/product"
	productService "demogo/internal/app/service/product"
	productRepository "demogo/internal/domain/product"
	productPb "demogo/internal/proto/product"
	"demogo/pkg/constant"
	ihttp "demogo/pkg/http"
	"demogo/pkg/logger"
	imongo "demogo/pkg/mongo"
	ipubsub "demogo/pkg/pubsub"
	iredis "demogo/pkg/redis"
	"demogo/pkg/safesql"
	"demogo/pkg/utils"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"net"
	"net/http"
	"time"
)

type GRPCServer struct {
	server              *grpc.Server
	subscriptionGateway ServiceSubscription
}

type GRPCServiceHandlerGroup struct {
	// Add all the service handlers here
}

// GracefulStop gracefully stop GRPC server
func (s *GRPCServer) GracefulStop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) Serve(port string) error {
	tcpListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s.subscriptionGateway.Subscribe()

	return s.server.Serve(tcpListener)
}

var (
	grpcMetrics = grpc_prometheus.NewServerMetrics()
)

const (
	healthCheckGrpcMethod string = "/grpc.health.v1.Health/Check"
)

// ServiceInjectGRPC is for injecting dependencies into each layer of GRPC service handler
func ServiceInjectGRPC(
	config *config.Config,
	dbObj map[string]safesql.IDatabase,
	redisdb map[string]iredis.IRedis,
	mongoDB imongo.IMongoDB,
) (GRPCServiceHandlerGroup, ServiceSubscription) {
	ctx := context.Background()

	pubsub := &ipubsub.Client{
		Topics:        make(map[string]*ipubsub.Topic),
		Subscriptions: make(map[string]*ipubsub.Subscription),
	}

	_, err := pubsub.NewClientWithRetry(ctx, config.ProjectID, config)
	if err != nil {
		logger.Log.Errorf("failed to create new client, error : %v", err)
	}

	_ = ihttp.HTTPCall{
		Conf: &config.HTTPConfig,
	}

	serviceHandlerGroup := GRPCServiceHandlerGroup{
		// Add all the service handlers here
	}

	gateway := NewSubGateway(pubsub, config)

	return serviceHandlerGroup, gateway
}

// InitServer for initializing server (initializes both GRPC and http server)
func InitServer(conf *config.Config, dbObj map[string]safesql.IDatabase,
	redisdb map[string]iredis.IRedis, mongoDB imongo.IMongoDB) (*GRPCServer, *http.Server) {

	_, pubsubGateway := ServiceInjectGRPC(conf, dbObj, redisdb, mongoDB)

	grpcServer := GRPCServer{
		server: grpc.NewServer(
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_prometheus.StreamServerInterceptor,
				grpctrace.StreamServerInterceptor(
					grpctrace.WithServiceName(conf.Core.Name),
					grpctrace.NonErrorCodes(codes.NotFound, codes.InvalidArgument, codes.FailedPrecondition),
				),
				grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(utils.GrpcPanicHandler())),
			)),
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				// grpc_logrus.UnaryServerInterceptor(logrusEntry, grpcLoggerOptions...),
				grpctrace.UnaryServerInterceptor(
					grpctrace.WithServiceName(conf.Core.Name),
					grpctrace.NonErrorCodes(codes.NotFound, codes.InvalidArgument, codes.FailedPrecondition),
				),
				UnaryRequestInterceptor(conf),
				grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(utils.GrpcPanicHandler())),
			)),
			grpc.KeepaliveParams(
				keepalive.ServerParameters{
					Time: 10000,
				},
			),
		),
		subscriptionGateway: pubsubGateway,
	}

	productRepo := productRepository.New(dbObj[constant.DatabasePostgreSQL], dbObj[constant.DatabasePostgreSQL])
	productSvc := productService.New(productRepo)
	productApi := productGrpc.New(conf, productSvc)
	productPb.RegisterProductServiceServer(grpcServer.server, &productApi)

	// healthpb.RegisterHealthServer(grpcServer.server, grpcServer)
	reflection.Register(grpcServer.server)

	// Initialize all metrics.
	grpcMetrics.EnableHandlingTimeHistogram()
	grpc_prometheus.EnableHandlingTimeHistogram()

	grpc_prometheus.Register(grpcServer.server)
	grpcMetrics.InitializeMetrics(grpcServer.server)

	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:              ":" + conf.Core.Port,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return &grpcServer, httpServer
}
