package http

import (
	"flag"
	"net/http"
	"time"
	"tools-home/internal/conf"
	"tools-home/internal/model"
	"tools-home/internal/server/http/common"
	"tools-home/internal/server/http/math/primary/grade1"
	"tools-home/internal/service"

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
	v1 := r.Group("/api/v1")
	v1.GET("/current-user", currentUser)
	{
		g := v1.Group("/time-converter")
		g.GET("", howToStart)
	}
	{
		g := v1.Group("/uuid-v4")
		g.GET("", generateUuid)
	}
	{
		g := v1.Group("/math")
		g.GET("/primary/grade1", grade1.Bit2Questions)
		g.GET("/primary/grade1/base", grade1.Bit2Questions)
	}
}

func ping(c *gin.Context) {
	if err := svc.Ping(); err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
	c.JSON(http.StatusOK, common.Reply{Data: "PONG"})
}

func howToStart(c *gin.Context) {
	k := &model.Bedrock{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(http.StatusOK, common.Reply{Data: k})
}
