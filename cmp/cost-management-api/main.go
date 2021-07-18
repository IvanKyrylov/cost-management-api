package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/IvanKyrylov/cost-management-api/internal/config"
	"github.com/IvanKyrylov/cost-management-api/internal/user"

	userdb "github.com/IvanKyrylov/cost-management-api/internal/user/db"
	"github.com/IvanKyrylov/cost-management-api/pkg/logging"
	"github.com/IvanKyrylov/cost-management-api/pkg/postgres"
	"github.com/IvanKyrylov/cost-management-api/pkg/shutdown"
)

// TEST DEV
// Test Home
func main() {
	logger := logging.Init()
	logging.CommonLog.Println("logger init")

	logging.CommonLog.Println("config init")
	cfg := config.GetConfig()

	logging.CommonLog.Println("router init")
	router := http.NewServeMux()

	pgClient, err := postgres.NewClient(context.Background(), cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SslMode)

	userStorage := userdb.NewStorage(mongoClient, cfg.MongoDB.CollectionUsers, logger)

	if err != nil {
		panic(err)
	}

	userService, err := user.NewService(userStorage, logger)

	if err != nil {
		panic(err)
	}

	userHandler := user.Handler{
		Logger:      logger,
		UserService: userService,
	}

	userHandler.Register(router)

	logger.Println("Start application")
	start(router, logger, cfg)
}

func start(router http.Handler, logger *log.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Printf("socket path: %s", socketPath)

		logger.Println("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		// logger.Printf("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)
		logger.Printf("bind application to host: %s and port: %s", "", os.Getenv("PORT"))

		var err error

		// listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", "", os.Getenv("PORT")))

		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Println("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
