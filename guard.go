package ginx

import (
	"github.com/gin-gonic/gin"
)

type Guard interface {
	Guard(c *gin.Context) (interface{}, E)
	Response(c *gin.Context, e error)
}

type GuardRes func(c *gin.Context, err error)

func ParamValidator(g Guard) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			param interface{}
			err   E
		)

		if param, err = g.Guard(c); !err.Empty() {
			g.Response(c, err)
			c.Abort()
			return
		}

		c.Set("param", param)
		c.Next()
	}
}

func G(g Guard) gin.HandlerFunc {
	return ParamValidator(g)
}
