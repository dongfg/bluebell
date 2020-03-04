package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server wrap of http.Server with gin
type Server struct {
	opts *Options
}

// Options server dependency
type Options struct {
	R    *gin.Engine
	Port int
}

// New server with dependency
func New(opts *Options) *Server {
	return &Server{
		opts: opts,
	}
}

// Start server gracefully
func (server *Server) Start() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", server.opts.Port),
		Handler: server.opts.R,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Start Server @ %s\n", srv.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("Server exiting")
}
