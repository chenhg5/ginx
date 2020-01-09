package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth interface {
	Check(c *gin.Context) bool
	User(c *gin.Context, userPointer interface{})
	Login(http *http.Request, w http.ResponseWriter, user interface{}) interface{}
	Logout(http *http.Request, w http.ResponseWriter) bool
	Middleware() gin.HandlerFunc
}

var defaultFilterRes = func(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  "请先登录",
	})
}

func RegisterJWTAuth(cfg JWTAuthConfig) Auth {
	return newJwtAuthDriver(cfg)
}
