package main

import (
	"context"
	"fmt"
	"github.com/dongfg/bluebell/config"
	"github.com/dongfg/bluebell/consul"
	"github.com/dongfg/bluebell/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var client *consul.Consul

func init() {
	client = consul.New(os.Getenv("CONSUL_ADDR"), os.Getenv("CONSUL_TOKEN"))
	if err := config.Load(client); err != nil {
		panic(err)
	}
}

func main() {
	r := setupRouter()
	controller.Register(r)
	setupServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", config.Basic.Port),
		Handler: r,
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFile("/", "public/index.html")
	r.StaticFile("/swagger.yml", "public/swagger.yml")
	return r
}

func setupServer(srv *http.Server) {
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	client.Register(consul.Service(config.Basic.Service))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Deregister("bluebell")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("Server exiting")
}
