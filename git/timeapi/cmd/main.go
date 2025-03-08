package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cheina97/timeserver/pkg/api"
	"github.com/cheina97/timeserver/pkg/flags"
	"github.com/cheina97/timeserver/pkg/handlers"
)

func main() {
	opts := flags.NewOptions()
	flags.Init(opts)

	router := gin.New()
	router.Use(gin.Logger())

	api.RegisterHandlers(router, handlers.NewServer())

	s := &http.Server{
		Handler:     router,
		Addr:        opts.Addr,
		ReadTimeout: opts.ReadTimeout,
	}

	log.Fatal(s.ListenAndServe())
}
