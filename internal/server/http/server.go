package http

import (
	"flag"
	"log"
	"net/http"
	"time"
	"tools-home/internal/conf"
	"tools-home/internal/model"
	"tools-home/internal/service"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	svc   *service.Service
	limit ratelimit.Limiter
	rps   int
)

func init() {
	flag.IntVar(&rps, "rps", 100, "request per second")
}

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}

// New new a bm server.
func New(s *service.Service) (engine *http.Server, err error) {
	var cfg struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	limit = ratelimit.New(rps)

	if err = conf.Load("http.json", &cfg); err != nil {
		return
	}
	svc = s

	router := gin.Default()
	router.Use(leakBucket())
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
	{
		g := r.Group("/tools-home")
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
