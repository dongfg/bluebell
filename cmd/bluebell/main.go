package main

import (
	"context"
	"fmt"
	"github.com/dongfg/bluebell/internal/config"
	"github.com/dongfg/bluebell/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if err := config.Load(); err != nil {
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
	box := packr.New("static", "../../static")
	r.GET("/", func(c *gin.Context) {
		src, _ := box.Find("index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", src)
	})
	r.GET("/swagger.yml", func(c *gin.Context) {
		src, _ := box.Find("swagger.yml")
		c.Data(http.StatusOK, "text/vnd.yaml; charset=utf-8", src)
	})
	return r
}

func setupServer(srv *http.Server) {
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
