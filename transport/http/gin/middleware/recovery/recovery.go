/**
 * @Time: 2021/2/25 4:32 下午
 * @Author: varluffy
 */

package recovery

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/rich/errcode"
	"github.com/varluffy/rich/log"
	"github.com/varluffy/rich/transport/http/gin/ginx"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"runtime"
)

type Option func(*options)

type options struct {
	logger *zap.Logger
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func Recovery(opts ...Option) gin.HandlerFunc {
	options := &options{
		logger: log.NewLogger(),
	}
	for _, o := range opts {
		o(options)
	}
	logger := options.logger
	logger = logger.With(zap.String("module", "middleware/recovery"))
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
				ginx.ErrorResponse(c, errcode.New(500, "inner error", http.StatusInternalServerError))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
