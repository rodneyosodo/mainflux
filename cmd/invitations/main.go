// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

// Package main contains invitations main function to start the invitations service.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"

	chclient "github.com/absmach/callhome/pkg/client"
	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/internal"
	"github.com/absmach/magistrala/internal/clients/jaeger"
	clientspg "github.com/absmach/magistrala/internal/clients/postgres"
	"github.com/absmach/magistrala/internal/postgres"
	"github.com/absmach/magistrala/internal/server"
	"github.com/absmach/magistrala/internal/server/http"
	"github.com/absmach/magistrala/invitations"
	"github.com/absmach/magistrala/invitations/api"
	"github.com/absmach/magistrala/invitations/middleware"
	invitationspg "github.com/absmach/magistrala/invitations/postgres"
	mglog "github.com/absmach/magistrala/logger"
	"github.com/absmach/magistrala/pkg/auth"
	mgsdk "github.com/absmach/magistrala/pkg/sdk/go"
	"github.com/absmach/magistrala/pkg/uuid"
	"github.com/caarlos0/env/v10"
	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/pyroscope-go"
	"github.com/jmoiron/sqlx"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "invitations"
	envPrefixDB    = "MG_INVITATIONS_DB_"
	envPrefixHTTP  = "MG_INVITATIONS_HTTP_"
	envPrefixAuth  = "MG_AUTH_GRPC_"
	defDB          = "invitations"
	defSvcHTTPPort = "9020"
)

type config struct {
	LogLevel      string  `env:"MG_INVITATIONS_LOG_LEVEL"      envDefault:"info"`
	UsersURL      string  `env:"MG_USERS_URL"                  envDefault:"http://localhost:9002"`
	DomainsURL    string  `env:"MG_DOMAINS_URL"                envDefault:"http://localhost:8189"`
	InstanceID    string  `env:"MG_INVITATIONS_INSTANCE_ID"    envDefault:""`
	JaegerURL     url.URL `env:"MG_JAEGER_URL"                 envDefault:"http://localhost:14268/api/traces"`
	TraceRatio    float64 `env:"MG_JAEGER_TRACE_RATIO"         envDefault:"1.0"`
	SendTelemetry bool    `env:"MG_SEND_TELEMETRY"             envDefault:"true"`
	LokiURL       string  `env:"GOPHERCON_LOKI_URL"            envDefault:""`
	PyroScopeURL  string  `env:"GOPHERCON_PYROSCOPE_URL"       envDefault:""`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err)
	}

	var level slog.Level
	err := level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		log.Fatalf("failed to parse log level: %s", err.Error())
	}
	fanout := slogmulti.Fanout(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)
	if cfg.LokiURL != "" {
		config, err := loki.NewDefaultConfig(cfg.LokiURL)
		if err != nil {
			log.Fatalf("failed to create loki config: %s", err.Error())
		}
		config.TenantID = svcName
		client, err := loki.New(config)
		if err != nil {
			log.Fatalf("failed to create loki client: %s", err.Error())
		}

		hander := slogloki.Option{Level: level, Client: client}.NewLokiHandler()
		fanout = slogmulti.Fanout(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: level,
			}),
			hander,
		)
	}

	logger := slog.New(fanout).With("service", svcName)
	slog.SetDefault(logger)

	var exitCode int
	defer mglog.ExitWithError(&exitCode)

	if cfg.InstanceID == "" {
		if cfg.InstanceID, err = uuid.New().ID(); err != nil {
			logger.Error(fmt.Sprintf("failed to generate instanceID: %s", err))
			exitCode = 1
			return
		}
	}

	dbConfig := clientspg.Config{Name: defDB}
	if err := env.ParseWithOptions(&dbConfig, env.Options{Prefix: envPrefixDB}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s database configuration : %s", svcName, err))
		exitCode = 1
		return
	}
	db, err := clientspg.Setup(dbConfig, *invitationspg.Migration())
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer db.Close()

	authConfig := auth.Config{}
	if err := env.ParseWithOptions(&authConfig, env.Options{Prefix: envPrefixAuth}); err != nil {
		logger.Error(fmt.Sprintf("failed to load auth configuration : %s", err.Error()))
		exitCode = 1
		return
	}
	authClient, authHandler, err := auth.Setup(authConfig)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer authHandler.Close()
	logger.Info("Successfully connected to auth grpc server " + authHandler.Secure())

	tp, err := jaeger.NewProvider(ctx, svcName, cfg.JaegerURL, cfg.InstanceID, cfg.TraceRatio)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init Jaeger: %s", err))
		exitCode = 1
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)

	if cfg.PyroScopeURL != "" {
		if _, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: svcName,
			ServerAddress:   cfg.PyroScopeURL,
			Logger:          nil,
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
			},
		}); err != nil {
			log.Fatalf("failed to start pyroscope: %s", err.Error())
		}
	}

	svc, err := newService(db, dbConfig, authClient, tracer, cfg, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create %s service: %s", svcName, err))
		exitCode = 1
		return
	}

	httpServerConfig := server.Config{Port: defSvcHTTPPort}
	if err := env.ParseWithOptions(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err))
		exitCode = 1
		return
	}

	httpSvr := http.New(ctx, cancel, svcName, httpServerConfig, api.MakeHandler(svc, logger, cfg.InstanceID), logger)

	if cfg.SendTelemetry {
		chc := chclient.New(svcName, magistrala.Version, logger, cancel)
		go chc.CallHome(ctx)
	}

	g.Go(func() error {
		return httpSvr.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, httpSvr)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("%s service terminated: %s", svcName, err))
	}
}

func newService(db *sqlx.DB, dbConfig clientspg.Config, authClient magistrala.AuthServiceClient, tracer trace.Tracer, conf config, logger *slog.Logger) (invitations.Service, error) {
	database := postgres.NewDatabase(db, dbConfig, tracer)
	repo := invitationspg.NewRepository(database)

	config := mgsdk.Config{
		UsersURL:   conf.UsersURL,
		DomainsURL: conf.DomainsURL,
	}
	sdk := mgsdk.NewSDK(config)

	svc := invitations.NewService(repo, authClient, sdk)
	svc = middleware.Tracing(svc, tracer)
	svc = middleware.Logging(logger, svc)
	counter, latency := internal.MakeMetrics(svcName, "api")
	svc = middleware.Metrics(counter, latency, svc)

	return svc, nil
}
