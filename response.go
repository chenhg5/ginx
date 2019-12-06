package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSON(c *gin.Context, data interface{}, msg string, code Code) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
	return
}

func OKWithData(c *gin.Context, data interface{}) {
	JSON(c, data, "ok", 0)
	return
}

func OKWithData200(c *gin.Context, data interface{}) {
	JSON(c, data, "ok", 200)
	return
}

func OK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
	})
	return
}

func Empty(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": []gin.H{},
		"msg":  "",
	})
	return
}

func Error(c *gin.Context, err E) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": err.Code(),
		"data": []gin.H{},
		"msg":  err.Error(),
	})
	return
}

func CommonError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"msg":  msg,
	})
	return
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  "参数错误",
		"data": []gin.H{},
	})
	return
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"msg":  "请先登录",
		"data": gin.H{},
		"code": 401,
	})
	return
}
