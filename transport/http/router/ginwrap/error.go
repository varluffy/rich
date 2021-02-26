/**
 * @Time: 2021/2/26 11:32 上午
 * @Author: varluffy
 * @Description: ginwarp
 */

package ginwrap

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/errcode"
	"net/http"
)

func Response(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": data})
}

func ErrorResponse(c *gin.Context, err error, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	if data == nil {
		data = gin.H{}
	}
	if e := new(errcode.Error); errors.As(err, &e) {
		c.JSON(e.Status, gin.H{"code":e.Code, "message":e.Message, "data":data})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code":500, "message":"unknown", "data":gin.H{}})
}

