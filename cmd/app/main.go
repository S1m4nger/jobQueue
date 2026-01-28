package main

import (
	"context"
	apphttp "job-queue/internal/http"
	"job-queue/internal/http/handlers"
	"job-queue/internal/queue"
	gormrepo "job-queue/internal/repository/gorm"
	"job-queue/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 1. Database setup (GORM + SQLite)
	gdb, err := gorm.Open(sqlite.Open("./jobs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// 2. Initialize layers
	repo := gormrepo.NewJobRepository(gdb)
	q := queue.NewMemoryQueue(100)
	svc := service.NewJobService(repo, q)
	handler := handlers.NewJobHandler(svc)
	router := apphttp.NewRouter(handler)

	// 3. Start workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc.StartWorkers(ctx, 3)

	// 4. Start Server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 5. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel() // Stop workers

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
