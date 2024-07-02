package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	driver "github.com/tjmaynes/shopping-cart-service-go/internal/driver"
	"github.com/tjmaynes/shopping-cart-service-go/internal/handler"
	handlers "github.com/tjmaynes/shopping-cart-service-go/internal/handler/http"
	"github.com/tjmaynes/shopping-cart-service-go/internal/pkg/item"
)

// API ..
type API struct {
	DbConn  *sql.DB
	Handler http.Handler
}

// NewAPI ..
func NewAPI(dbSource string) *API {
	dbConn, err := driver.ConnectDB(dbSource)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	cartRepository := item.NewRepository(dbConn)
	cartService := item.NewService(cartRepository)
	cartHandler := handlers.NewItemHandler(cartService)

	healthCheckHandler := handlers.NewHealthCheckHandler(dbConn)

	return &API{
		DbConn:  dbConn,
		Handler: handler.Initialize(cartHandler, healthCheckHandler),
	}
}

// Run ..
func (a *API) Run(serverPort string) {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", serverPort),
		Handler:        a.Handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println(fmt.Sprintf("Running server on port %s...", serverPort))

	idleConnsClosed := make(chan struct{})
	go setupGracefulShutdown(server, a.DbConn, idleConnsClosed)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Server closed: %v", err)
	}

	<-idleConnsClosed
}

func setupGracefulShutdown(server *http.Server, db *sql.DB, idleConnsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	<-sigint

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	defer db.Close()

	close(idleConnsClosed)
}
