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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {

	// Load application configuration
	cfg := config.Load()

	ctx := context.Background()

	db, err := repo.NewDb(ctx, cfg)
	if err != nil {
		fmt.Printf("error in db connection: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	awsSession, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String(cfg.AWSRegion),
				Credentials: credentials.NewStaticCredentials(
					cfg.AWSAccessKey,
					cfg.AWSSecretKey,
					"",
				),
			},
		},
	)

	if err != nil {
		fmt.Printf("Error creating AWS session: %v\n", err)
		os.Exit(1)
	}

	uploader := s3manager.NewUploader(awsSession)

	awsRepo := repo.NewAwsProjRepo(db)

	awsHandler := handler.NewAwsProjHandler(
		awsRepo,
		uploader,
		cfg.AWSBucketName,
		cfg.AWSRegion,
	)

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

		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			fmt.Println("ListenAndServe error:", err)
		}

	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-quit

	fmt.Println("Received signal:", sig)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.JWTTokenDuration,
	)

	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting:", sig)
}
// 