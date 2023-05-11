package app

import (
	"context"
	"github.com/Arkosh744/chat-server/internal/closer"
	config "github.com/Arkosh744/chat-server/internal/config"
	"go.uber.org/zap"
	"log"
)

type App struct {
	serviceProvider *serviceProvider
	log             *zap.SugaredLogger
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return nil
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		config.Init,
		app.initLogger,
		app.initServiceProvider,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initLogger(_ context.Context) error {
	zapLog, err := config.SelectLogger()
	if err != nil {
		log.Fatalf("failed to get logger: %s", err.Error())
	}

	app.log = zapLog.Sugar()

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider(app.log)

	return nil
}
