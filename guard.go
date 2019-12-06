package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Guard interface {
	Guard(c *gin.Context) (interface{}, error)
}

type GuardRes func(c *gin.Context, err error)

var defaultGuardRes = func(c *gin.Context, err error) {
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
			err   error
		)

		if param, err = g.Guard(c); err != nil {
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
