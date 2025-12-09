package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kralle333/keyvaluestore/internal/handler"
	"github.com/kralle333/keyvaluestore/internal/model"
	"github.com/kralle333/keyvaluestore/internal/repository"
	"github.com/kralle333/keyvaluestore/internal/service"
	"go.uber.org/zap"
)

type App struct {
	Config         AppConfig
	Logger         *zap.Logger
	Communication  *model.KeyValueActorCommunication
	Storage        *repository.KeyValueStorage
	Actor          *service.KeyValueActor
	SnapshotLogger *service.SnapshotService
	HttpServer     *handler.KeyValueHttpServer
}

func NewApp(config AppConfig) (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Setup dependencies
	communication := model.NewKeyValueActorCommunication()
	storage := repository.NewKeyValueStorage(config.SnapshotDir, *logger)

	// Inject dependencies and use config values
	server := handler.NewKeyValueHttpServer(communication, config.ListeningPort)
	snapshotLogger := service.NewSnapshotService(communication, config.SnapshotIntervalSeconds, logger)
	actor := service.NewKeyValueActor(communication, storage, logger)

	// If we already have saved snapshot, use it to populate Actor's internal data
	if config.RestoreFromSnapshot {
		snapshot, err := storage.RetrieveLatest()
		if err != nil {
			return nil, err
		}
		actor.PopulateFromSnapshot(snapshot)
	}

	return &App{
		Config:         config,
		Logger:         logger,
		Storage:        storage,
		Actor:          actor,
		SnapshotLogger: snapshotLogger,
		HttpServer:     server,
	}, nil
}

func (app *App) Run() error {
	app.Logger.Info("Starting application services")

	app.Actor.Spawn()
	app.SnapshotLogger.Spawn()

	go func() {
		if err := app.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("HTTP server failed to listen and serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutdown signal received, starting graceful shutdown...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown select loops
	app.Actor.Shutdown()
	app.SnapshotLogger.Shutdown()

	app.Logger.Info("Application gracefully stopped")
	app.Logger.Sync()
	return nil
}
