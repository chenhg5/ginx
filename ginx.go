package ginx

import (
	"context"
	"github.com/gavv/httpexpect"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"testing"
	"time"
)

type GinX struct {
	Router *gin.Engine

	config Config
}

type process func(*gin.Engine)

func New() *GinX {
	return &GinX{
		Router: gin.New(),
	}
}

func (g *GinX) InjectRouter(proc process) *GinX {
	proc(g.Router)
	return g
}

func (g *GinX) SetConfig(config Config) *GinX {
	g.config = config
	return g
}

func (g *GinX) Use(middleware ...gin.HandlerFunc) *GinX {
	g.Router.Use(middleware...)
	return g
}

func (g *GinX) Run(port ...string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if g.config.Debug() {
		gin.SetMode(gin.DebugMode)
		g.Router.Use(gin.Logger())
		pprof.Register(g.Router)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if len(port) == 0 {
		port = []string{"80"}
	}

	srv := &http.Server{
		Addr:    ":" + port[0],
		Handler: g.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func NewWithRouter(router *gin.Engine) *GinX {
	return &GinX{
		Router: router,
	}
}

func (g *GinX) Expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(g.Router),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}
