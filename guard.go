package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Guard interface {
	Guard(c *gin.Context) (interface{}, E)
}

type GuardRes func(c *gin.Context, err E)

var defaultGuardRes = func(c *gin.Context, err E) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func SetDefaultGuardRes(d GuardRes) {
	defaultGuardRes = d
}

func ParamValidator(g Guard) gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			param interface{}
			err   E
		)

		if param, err = g.Guard(c); !err.Empty() {
			defaultGuardRes(c, err)
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
