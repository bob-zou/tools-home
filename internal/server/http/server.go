package http

import (
	"net/http"
	"time"
	"tools-home/internal/conf"
	"tools-home/internal/model"
	"tools-home/internal/service"

	"github.com/gin-gonic/gin"
)

var svc *service.Service

// New new a bm server.
func New(s *service.Service) (engine *http.Server, err error) {
	var cfg struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	if err = conf.Load("http.json", &cfg); err != nil {
		return
	}
	svc = s

	router := gin.Default()
	initRouter(router)
	engine = &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	}

	go func() {
		if err = engine.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return
}

func initRouter(r *gin.Engine) {
	r.GET("/monitor/ping", ping)
	g := r.Group("/tools-home")
	{
		g.GET("/start", howToStart)
	}
}

func ping(c *gin.Context) {
	if err := svc.Ping(); err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
	c.JSON(http.StatusOK, CommonResponse{Data: "PONG"})
}

func howToStart(c *gin.Context) {
	k := &model.Bedrock{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(http.StatusOK, CommonResponse{Data: k})
}