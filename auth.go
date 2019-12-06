package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Auth interface {
	Check(c *gin.Context) bool
	User(c *gin.Context) interface{}
	Login(http *http.Request, w http.ResponseWriter, user map[string]interface{}) interface{}
	Logout(http *http.Request, w http.ResponseWriter) bool
}

var jwtAuth Auth

func RegisterJWTAuth(secret, alg, header string, exp time.Duration) {
	jwtAuth = newJwtAuthDriver(secret, alg, header, exp)
}

var defaultFilterRes = func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "请先登录",
	})
}

func SetDefaultFilterRes(d gin.HandlerFunc) {
	defaultFilterRes = d
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !jwtAuth.Check(c) {
			defaultFilterRes(c)
			c.Abort()
		}
		c.Next()
	}
}

func User(c *gin.Context) map[string]interface{} {
	return jwtAuth.User(c).(map[string]interface{})
}
