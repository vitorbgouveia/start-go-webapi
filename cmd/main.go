package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vitorbgouveia/start-project-go/internal/repositories"
	pkg "github.com/vitorbgouveia/start-project-go/package"
	"github.com/vitorbgouveia/start-project-go/package/routes"
	"github.com/vitorbgouveia/start-project-go/package/services"
)

var (
	dbUser string
	dbPass string
	dbName string
	dbPort int
	dbHost string
)

func init() {
	flag.StringVar(&dbUser, "db_user", "root", "user to connect database")
	flag.StringVar(&dbPass, "db_pass", "root", "password to connect database")
	flag.StringVar(&dbName, "db_name", "start_project", "database name")
	flag.StringVar(&dbHost, "db_host", "localhost", "port to connect database")
	flag.IntVar(&dbPort, "db_port", 5432, "port to connect database")
	flag.Parse()
}

func main() {
	logger := pkg.NewLogger()
	defer logger.Sync()

	ctx := context.Background()
	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalw("could not connect database", zap.Error(err))
	}

	dbPoolConn, err := db.DB()
	if err != nil {
		logger.Fatalw("could not create pool connections")
	}

	dbPoolConn.SetConnMaxIdleTime(time.Second * 3)
	dbPoolConn.SetMaxOpenConns(10)
	dbPoolConn.SetMaxIdleConns(5)
	defer dbPoolConn.Close()

	walletRepository := repositories.NewWalletRepository(dbPoolConn)

	walletService := services.NewWalletServie(logger, walletRepository)

	wsContainer := restful.NewContainer()
	routes.DeclareAllRoutes(logger, wsContainer, walletService)

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		Container:      wsContainer,
	}
	wsContainer.Filter(cors.Filter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	logger.Info("listening port 5000...")

	server := &http.Server{
		Addr:    ":5000",
		Handler: wsContainer,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-sigCtx.Done()
	logger.Info("got interruption signal")
	if err := server.Shutdown(ctx); err != nil {
		logger.Warn("server shutdown returned an err", zap.Error(err))
	}
	logger.Info("finishing app...")
}
