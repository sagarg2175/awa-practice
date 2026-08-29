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
	"time"
)

func init() {
	config.Load()
}

func main() {
	ctx := context.Background()
	db, err := repo.NewDb(ctx)
	if err != nil {
		fmt.Errorf("error in db connection %s", err)
	}

	awsRepo := repo.NewAwsProjRepo(db)
	awsHandler := handler.NewAwsProjHandler(awsRepo)

	router, err := handler.NewRouter(*awsHandler)
	if err != nil {
		fmt.Errorf("Error initializing router %s", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	duration, err := time.ParseDuration(os.Getenv("SHUTDOWN_TIME"))
	if err != nil {
		fmt.Println("Error in parsing duration", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	<-ctx.Done()
	fmt.Println("Server exiting", sig)

	defer db.Close()
}
