package main

import (
	"context"
	"database/sql"
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
	dbTxhistory "github.com/IvanKyrylov/cost-management-api/internal/user/db/txhistory"
	dbUser "github.com/IvanKyrylov/cost-management-api/internal/user/db/user"
	dbWallet "github.com/IvanKyrylov/cost-management-api/internal/user/db/wallet"

	"github.com/IvanKyrylov/cost-management-api/pkg/logging"
	"github.com/IvanKyrylov/cost-management-api/pkg/postgres"
	"github.com/IvanKyrylov/cost-management-api/pkg/shutdown"
)

func main() {
	logger := logging.Init()
	logging.CommonLog.Println("logger init")

	logging.CommonLog.Println("config init")
	cfg := config.GetConfig()

	logging.CommonLog.Println("router init")
	router := http.NewServeMux()

	pgClient, err := postgres.NewClient(context.Background(), cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SslMode)

	var userStorage user.UserStorage = dbUser.NewStorage(pgClient, logger)
	var walletStorage user.WalletStorage = dbWallet.NewStorage(pgClient, logger)
	var txhistoryStorage user.TransactionHistoryStorage = dbTxhistory.NewStorage(pgClient, logger)

	if err != nil {
		panic(err)
	}

	userService, err := user.NewService(userStorage, walletStorage, txhistoryStorage, logger)

	if err != nil {
		panic(err)
	}

	userHandler := user.Handler{
		Logger:      logger,
		UserService: userService,
	}

	userHandler.Register(router)

	logger.Println("Start application")
	start(router, logger, cfg, pgClient)
}

func start(router http.Handler, logger *log.Logger, cfg *config.Config, client *sql.DB) {
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
		logger.Printf("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			client.Close()
			logger.Println("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
