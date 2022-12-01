package main

import (
	"context"
	"github.com/BlackRRR/transaction-service/internal/config"
	"github.com/BlackRRR/transaction-service/internal/log"
	"github.com/BlackRRR/transaction-service/internal/repository"
	"github.com/BlackRRR/transaction-service/internal/server"
	"github.com/BlackRRR/transaction-service/internal/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//new logger
	logger := log.NewProductionLogger(nil)

	//init config
	initConfig, err := config.InitConfig()
	if err != nil {
		logger.Sugar().Fatalf("failed to init config: %s", err.Error())
	}

	//init db
	DBPool, err := repository.InitDB(initConfig.PGConn)
	if err != nil {
		logger.Sugar().Fatalf("failed to init db: %s", err.Error())
	}

	//init repository
	repo, err := repository.NewRepo(context.Background(), DBPool)
	if err != nil {
		logger.Sugar().Fatalf("failed to init transaction repo %s", err.Error())
	}

	//init service
	reader := services.NewReader(logger, repo)
	reader.ReadQueue()

	err = reader.RecoveryUncompletedTransactions()
	if err != nil {
		logger.Sugar().Infof("failed to recover transactions %s", err.Error())
	}

	//init new controller
	controller := server.NewController(reader)

	//init http server
	httpServer := server.New(controller.InitRoutes(), server.Port(initConfig.ServicePort))
	logger.Sugar().Infof("server started on %s port", initConfig.ServicePort)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	//waiting for server shutdown or system signal
	select {
	case s := <-interrupt:
		logger.Sugar().Warnf("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		logger.Sugar().Warnf("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Sugar().Warnf("app - Run - httpServer.Shutdown: %v", err)
		return
	}
	logger.Sugar().Info("server shutdown")
}
