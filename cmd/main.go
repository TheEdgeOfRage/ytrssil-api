package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

func init() {
	// always use UTC
	time.Local = time.UTC
}

func fetcherRoutine(l log.Logger, h handler.Handler) {
	for {
		err := h.FetchVideos(context.Background())
		if err != nil {
			l.Log("level", "ERROR", "function", "main.fetcherRoutine", "call", "handler.FetchVideos", "err", err)
		}
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	log := log.NewLogger()

	config, err := config.Parse()
	if err != nil {
		log.Log("level", "FATAL", "call", "config.Parse", "error", err)
		return
	}
	db, err := db.NewPostgresDB(log, config.DB)
	if err != nil {
		log.Log("level", "FATAL", "call", "db.NewPostgresDB", "error", err)
		return
	}
	parser := feedparser.NewParser(log)
	handler := handler.New(log, db, parser)
	gin.SetMode(gin.ReleaseMode)
	router, err := ytrssil.SetupGinRouter(
		log,
		handler,
		auth.AuthMiddleware(db),
	)
	if err != nil {
		log.Log("level", "FATAL", "call", "ytrssil.SetupGinServer", "error", err)
		return
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Gin.Port),
		Handler: router,
	}
	server.RegisterOnShutdown(func() {
		log.Log("level", "INFO", "msg", "shutdown server.Close()")
	})

	quit := make(chan os.Signal, 1)
	// handle Interrupt (ctrl-c) Term, used by `kill` et al, HUP which is commonly used to reload configs
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		s := <-quit
		log.Log(
			"level", "INFO",
			"msg", "signalRecv, for quitting",
			"signal", s,
		)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Log(
				"level", "ERROR",
				"call", "server.Shutdown",
				"error", err,
			)
		}
	}()

	// start periodic fetch videos routine
	go fetcherRoutine(log, handler)

	log.Log(
		"level", "INFO",
		"msg", "ytrssil API is starting up",
		"port", config.Gin.Port,
	)
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Log(
				"level", "ERROR",
				"call", "server.ListenAndServe",
				"error", err,
			)
		}
	}
	log.Log("level", "INFO", "msg", "exit complete")
}
