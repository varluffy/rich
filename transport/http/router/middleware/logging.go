/**
 * @Time: 2021/2/25 2:56 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/varluffy/ginx/log"
	"go.uber.org/zap"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger = logger.WithOptions(zap.WithCaller(false))
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter
		traceId := c.GetHeader("X-Request-Id")
		if traceId == "" {
			traceId = uuid.NewV4().String()
		}
		start := time.Now()
		ctx := log.NewTraceIDContext(c.Request.Context(), traceId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		p := c.Request.URL.Path
		method := c.Request.Method
		body, _ := c.GetRawData()
		msg := fmt.Sprintf("[HTTP] %s-%s-%s-%d (%dms)", p, method, c.ClientIP(), c.Writer.Status(), time.Since(start)/1e6)
		fields := []zap.Field{
			zap.String("trace_id", traceId),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Any("header", c.Request.Header),
			zap.String("request", c.Request.PostForm.Encode()),
			zap.String("body", string(body)),
			zap.String("response", bodyWriter.body.String()),
			zap.String("ip", c.ClientIP()),
		}
		if c.Writer.Status() <= 200 {
			logger.Info(msg, fields...)
		} else if c.Writer.Status() < 500 {
			logger.Warn(msg, fields...)
		} else {
			logger.Error(msg, fields...)
		}
	}
}
