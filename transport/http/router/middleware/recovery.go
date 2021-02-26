/**
 * @Time: 2021/2/25 4:32 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/errcode"
	"github.com/varluffy/ginx/log"
	"github.com/varluffy/ginx/transport/http/router/ginwrap"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"runtime"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			ctx := c.Request.Context()
			if err := recover(); err != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Error("[Recovery from panic]",
					zap.String("trace_id", log.FromTraceIDContext(ctx)),
					zap.String("request", string(httpRequest)),
					zap.Any("err", err),
					zap.String("stack", string(buf)),
				)
				ginwrap.ErrorResponse(c, errcode.New(500, "inner error", http.StatusInternalServerError))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}