package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"tutorgpt/pkg/handler"
	"tutorgpt/pkg/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	utils.PrintConfig(config)

	chatHandler := handler.NewChatHandler(config)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/chat", chatHandler.HandleChat)

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	frontendPaths := []string{
		filepath.Join(exeDir, "frontend"),
		filepath.Join(exeDir, "..", "frontend"),
	}

	var frontendPath string
	for _, path := range frontendPaths {
		if _, err := os.Stat(path); err == nil {
			frontendPath = path
			break
		}
	}

	if frontendPath == "" {
		log.Fatalf("Could not find frontend directory. Tried paths: %v", frontendPaths)
	}

	log.Printf("Serving frontend from: %s", frontendPath)

	fs := http.FileServer(http.Dir(frontendPath))
	mux.Handle("/", fs)

	corsHandler := utils.CorsMiddleware(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: corsHandler,
	}

	go func() {
		log.Printf("Server starting on port %s...\n", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v\n", config.Port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")
}
