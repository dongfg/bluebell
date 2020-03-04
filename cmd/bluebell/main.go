package main

import (
	"github.com/dongfg/bluebell/internal/server"
	"log"
)

func main() {
	conf, err := initConfig()
	if err != nil {
		log.Panicln(err)
	}

	r := initRouter(conf)

	ctrl, err := initController(conf, r)
	if err != nil {
		log.Panicln(err)
	}
	ctrl.Register()

	srv := server.New(&server.Options{
		R:    r,
		Port: conf.Port,
	})
	srv.Start()
}
