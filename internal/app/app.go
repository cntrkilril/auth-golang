package app

import (
	"context"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	v1 "github/cntrkilril/auth-golang/internal/controller/http/v1"
	"github/cntrkilril/auth-golang/internal/infrastructure"
	"github/cntrkilril/auth-golang/internal/service"
	"github/cntrkilril/auth-golang/pkg/govalidator"
	"github/cntrkilril/auth-golang/pkg/hasher"
	"github/cntrkilril/auth-golang/pkg/tokens"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	// logger
	atom := zap.NewAtomicLevel()
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		atom,
	)
	logger := zap.New(zapCore)
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	l := logger.Sugar()
	atom.SetLevel(zapcore.Level(*cfg.Logger.Level))

	l.Infof("logger initialized successfully")

	// fiber
	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          v1.HandleError,
	})
	f.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))
	f.Use(cors.New(cors.Config{
		AllowHeaders: "*",
	}))

	l.Infof("fiber initialized successfully")

	// validator
	val := govalidator.New()
	l.Infof("validator initialized successfully")

	// hasher
	h := hasher.New(cfg.Hasher.Cost)
	l.Infof("hasher initialized successfully")

	// tokensWorker
	tokensWorker := tokens.New(cfg.Tokens.JWTKey, h)

	// mongo
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.Mongo.ConnString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		l.Fatal(err.Error())
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			l.Fatal(err.Error())
		}
	}()
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		l.Fatal(err.Error())
	}
	l.Infof("mongo connected successfully")

	//infrastructures
	tokenRepo := infrastructure.NewTokenRepository(*client.Database("auth"))

	// services
	tokenService := service.NewTokenService(tokenRepo, cfg.Tokens.ExpiresInAccessToken, cfg.Tokens.ExpiresInRefreshToken, h, tokensWorker)

	// controllers
	tokenHandler := v1.NewTokenHandler(tokenService, val)

	// groups
	apiGroup := f.Group("api")
	tokenGroup := apiGroup.Group("tokens")

	tokenHandler.Register(tokenGroup)

	go func() {
		err = f.Listen(net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port))
		if err != nil {
			l.Fatal(err.Error())
		}
	}()

	l.Debug("Started HTTP server")

	l.Debug("Application has started")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("Application has been shut down")

}
