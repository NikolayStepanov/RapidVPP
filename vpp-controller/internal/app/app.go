package app

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NikolayStepanov/RapidVPP/internal/config"
	"github.com/NikolayStepanov/RapidVPP/internal/delivery/http/handlers"
	"github.com/NikolayStepanov/RapidVPP/internal/infrastructure/vpp"
	"github.com/NikolayStepanov/RapidVPP/internal/mw"
	"github.com/NikolayStepanov/RapidVPP/internal/server"
	"github.com/NikolayStepanov/RapidVPP/internal/service"
	"github.com/NikolayStepanov/RapidVPP/internal/service/vpp/info"
	"github.com/NikolayStepanov/RapidVPP/internal/service/vpp/interfaces"
	ipServ "github.com/NikolayStepanov/RapidVPP/internal/service/vpp/ip"
	"github.com/NikolayStepanov/RapidVPP/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type App struct {
	config    *config.Config
	server    *server.Server
	services  *service.Services
	handler   *handlers.Handler
	vppClient *vpp.Client
}

func NewApp(config *config.Config) (*App, error) {
	VPPClient, err := vpp.NewClient(config.VPP.Socket)
	if err != nil {
		log.Fatalf("failed to create VPP client: %v", err)
	}

	infoService := info.NewService(VPPClient)
	interfaceService := interfaces.NewService(VPPClient)
	IPService := ipServ.NewService(VPPClient)
	handler := handlers.NewHandler(infoService, interfaceService, IPService)
	server := server.NewServer(config, mw.LoggerMiddleware(handler))
	return &App{
		config:    config,
		server:    server,
		handler:   handler,
		vppClient: VPPClient,
	}, nil
}

func Run() {
	cfg := config.Init()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger.InitLogger(&cfg.Logger)
	defer logger.Sync()
	watcherLoggerFile := logger.ObserveLoggerConfigFile(ctx, cfg)
	defer func(watcherLoggerFile *fsnotify.Watcher) {
		err := watcherLoggerFile.Close()
		if err != nil {
			panic(err)
		}
	}(watcherLoggerFile)

	app, err := NewApp(cfg)
	if err != nil {
		logger.Fatal("error new app", zap.Error(err))
	}
	defer app.vppClient.Close()

	go func() {
		defer cancel()
		if err := app.server.Run(); err != nil {
			logger.Error("error occurred while running http server", zap.Error(err))
		}
	}()
	logger.Info("RapidVPP is running")
	logger.Info("server started")
	<-ctx.Done()
	logger.Info("shutting down server")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()
	wgShutdown := sync.WaitGroup{}
	wgShutdown.Add(1)
	go func() {
		defer wgShutdown.Done()
		if err = app.server.Stop(ctxShutdown); err != nil {
			logger.Error("error occurred on server shutting down", zap.Error(err))
		}
	}()
	wgShutdown.Wait()
	logger.Info("RapidVPP stopped")
}
