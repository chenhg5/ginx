package ginx

import (
	"github.com/gavv/httpexpect"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"testing"
)

type GinX struct {
	Router *gin.Engine

	config Config
}

func New() *GinX {
	return &GinX{
		Router: gin.New(),
	}
}

func (g *GinX) SetConfig(config Config) *GinX {
	g.config = config
	return g
}

func (g *GinX) Use(middleware ...gin.HandlerFunc) *GinX {
	g.Router.Use(middleware...)
	return g
}

func (g *GinX) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if g.config.Debug() {
		gin.SetMode(gin.DebugMode)
		g.Router.Use(gin.Logger())
		pprof.Register(g.Router)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
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
