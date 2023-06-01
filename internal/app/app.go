package app

import (
	"context"
	"net"
	"sync"

	"github.com/Arkosh744/chat-server/internal/closer"
	"github.com/Arkosh744/chat-server/internal/config"
	"github.com/Arkosh744/chat-server/internal/interceptor"
	"github.com/Arkosh744/chat-server/internal/log"
	"github.com/Arkosh744/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := app.RunGrpcServer()
		if err != nil {
			log.Fatalf("failed to run grpc server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		config.Init,
		log.InitLogger,
		app.initServiceProvider,
		app.initGRPCServer,
		app.initGRPCClient,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()

	return nil
}

func (app *App) initGRPCClient(ctx context.Context) error {
	app.serviceProvider.GetChatService(ctx)

	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	authInterceptor := interceptor.NewAuthInterceptor(app.serviceProvider.GetAuthClient(ctx))

	app.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	reflection.Register(app.grpcServer)

	chat_v1.RegisterChatV1Server(app.grpcServer, app.serviceProvider.GetChatImpl(ctx))

	return nil
}

func (app *App) RunGrpcServer() error {
	log.Infof("GRPC server listening on %s", app.serviceProvider.GetGRPCConfig().GetHost())

	list, err := net.Listen("tcp", app.serviceProvider.GetGRPCConfig().GetHost())
	if err != nil {
		return err
	}

	err = app.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
