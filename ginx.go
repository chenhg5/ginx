package ginx

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"runtime"
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
