package main

import (
	"context"
	"fmt"
	"github.com/dongfg/bluebell/consul"
	"github.com/dongfg/bluebell/controller"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	Port    int
	Service struct {
		Name          string
		Address       string
		Port          int
		CheckUrl      string `yaml:"check-url"`
		CheckInterval string `yaml:"check-interval"`
	}
}

var client *consul.Consul
var c config

func init() {
	client = consul.New(os.Getenv("CONSUL_ADDR"), os.Getenv("CONSUL_TOKEN"))
	c = loadConfig(client)
}

func main() {
	r := setupRouter()
	controller.Setup(r)
	setupServer(&http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: r,
	})
}

func loadConfig(client *consul.Consul) config {
	rawConfig := client.Fetch("config/bluebell/yaml")
	config := config{}
	err := yaml.Unmarshal([]byte(rawConfig), &config)
	if err != nil {
		panic(err)
	}
	return config
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

	client.Register(consul.Service(c.Service))
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
