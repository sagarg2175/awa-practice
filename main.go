package main

import (
	"context"
	"fmt"
	"main/config"
	"main/handler"
	"main/repo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	db, err := repo.NewDb(ctx, cfg)
	if err != nil {
		fmt.Printf("error in db connection: %v\n", err)
		os.Exit(1)
	}

	awsRepo := repo.NewAwsProjRepo(db)
	awsHandler := handler.NewAwsProjHandler(awsRepo)

	router, err := handler.NewRouter(*awsHandler)
	if err != nil {
		fmt.Printf("Error initializing router: %v\n", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: router,
	}

	fmt.Println("Starting server on", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("ListenAndServe error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	duration := cfg.JWTTokenDuration

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	<-ctx.Done()
	fmt.Println("Server exiting", sig)

	db.Close()
}
