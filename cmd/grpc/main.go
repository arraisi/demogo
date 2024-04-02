package main

import (
	"context"
	"github.com/arraisi/demogo/config"
	mongoSvc "github.com/arraisi/demogo/internal/app/service/mongo"
	"github.com/arraisi/demogo/internal/server"
	"github.com/arraisi/demogo/pkg/constant"
	"github.com/arraisi/demogo/pkg/logger"
	imongo "github.com/arraisi/demogo/pkg/mongo"
	"github.com/arraisi/demogo/pkg/observer/tracer"
	iredis "github.com/arraisi/demogo/pkg/redis"
	"github.com/arraisi/demogo/pkg/safesql"
	"github.com/arraisi/demogo/pkg/utils"
	"github.com/arraisi/demogo/pkg/utils/profiler"
	"github.com/arraisi/demogo/pkg/utils/queryutils"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeWaiting = 5 * time.Second
)

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("init config failed: %v", err)
	}

	// init ddog APM tracer
	if conf.Core.Environment != "local" {
		tracer.Init(conf.Datadog.APMHost, conf.Core.Name)
		defer tracer.Stop()
	}

	queryutils.New(conf)

	//init logging
	customLog, err := logger.InitLogger(conf.Core.LogLevel)
	if err != nil {
		log.Fatalf("init logging failed: %v", err)
	}
	logger.SetLogger(customLog)

	profiler.StartGoProfiler(conf)
	profiler.StartGoCloudProfiler(conf)

	dbList := make(map[string]safesql.IDatabase)
	dbPostgres := safesql.PostgreSQLHandler{}
	dbPostgres.ConnectDB(&conf.Core.DBPostgres.READ, &conf.Core.DBPostgres.WRITE)

	dbList[constant.DatabasePostgreSQL] = &dbPostgres

	redisDb := iredis.RedisHandler{
		Client: &redis.Client{},
	}
	redisDb.ConnectRedis(&conf.Core.Redis)

	redisList := make(map[string]iredis.IRedis)
	redisList[constant.RedisDefault] = &redisDb

	registeredCollections := mongoSvc.GetRegisteredCollection()

	mongoDB := imongo.DBHandler{
		Client: &mongo.Client{},
	}

	err = mongoDB.Connect(&conf.Core.Mongo, registeredCollections)
	if err != nil {
		log.Fatalf("init mongoDB failed: %v", err)
	}

	var iMongo imongo.IMongoDB = &mongoDB

	grpcServer, httpServer := server.InitServer(conf, dbList, redisList, iMongo)

	log.Printf("%s starting on GRPC Port %s", conf.Core.Name, conf.Core.GRPC.Port)
	log.Printf("%s starting on HTTP Port %s", conf.Core.Name, conf.Core.Port)
	ShutdownApp(grpcServer, httpServer, conf.Core.Name, conf.Core.GRPC.Port, dbList, redisList[constant.RedisDefault], &mongoDB)
}

// ShutdownApp for shutting down server gracefully
func ShutdownApp(grpcServer *server.GRPCServer, httpServer *http.Server, serverName, grpcPort string, db map[string]safesql.IDatabase, redisdb iredis.IRedis, mongoDB imongo.IMongoDB) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(grpcPort); err != nil {
			log.Fatalf("[ERROR][grpcServer] Serve: %s\n", err)
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR][httpServer] ListenAndServe: %s\n", err)
		}
	}()

	doneVal := <-done
	log.Printf("signal received: %v", doneVal)
	log.Printf("signal received: %s", utils.JSONString(doneVal))
	log.Printf("Stopping %s\n", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), timeWaiting)
	defer func() {
		// extra handling here
		for _, v := range db {
			v.Close()
		}
		redisdb.Close()
		mongoDB.Disconnect(context.Background())
		cancel()
	}()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to stop http service %s %+v\n", serverName, err)
	}

	log.Printf("%s http stopped successfully\n", serverName)
	log.Printf("%s grpc stopped successfully\n", serverName)
}
