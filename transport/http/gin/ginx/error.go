/**
 * @Time: 2021/2/26 11:32 上午
 * @Author: varluffy
 */

package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/varluffy/rich/errcode"
	"net/http"
)

func Response(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": data})
}

func ErrorResponse(c *gin.Context, err error, data ...interface{}) {
	var d interface{}
	d = gin.H{}
	if len(data) > 0 {
		d = data
	}
	_ = c.Error(err)
	if e := new(errcode.Error); errors.As(err, &e) {
		c.JSON(e.Status, gin.H{"code": e.Code, "message": e.Message, "data": d})
		c.Abort()
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "unknown", "data": gin.H{}})
	c.Abort()
}
