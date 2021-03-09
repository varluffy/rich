/**
 * @Time: 2021/2/25 2:56 下午
 * @Author: varluffy
 */

package logging

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/varluffy/rich/log"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Option func(*options)

type options struct {
	logger   *zap.Logger
	skipPath []string
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithSkipPath(paths []string) Option {
	return func(o *options) {
		o.skipPath = paths
	}
}

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

func Server(opts ...Option) gin.HandlerFunc {
	options := &options{
		logger:   log.NewLogger(),
		skipPath: []string{},
	}
	for _, opt := range opts {
		opt(options)
	}
	logger := options.logger.WithOptions(zap.WithCaller(false))
	logger = logger.With(zap.String("module", "middleware/logging"))
	return func(c *gin.Context) {
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

		path := c.Request.URL.Path
		method := c.Request.Method
		body, _ := c.GetRawData()
		msg := fmt.Sprintf("[HTTP] %s-%s-%s-%d (%dms)", path, method, c.ClientIP(), c.Writer.Status(), time.Since(start)/1e6)
		fields := []zap.Field{
			zap.String("trace_id", traceId),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Any("header", c.Request.Header),
			zap.String("request", c.Request.PostForm.Encode()),
			zap.String("body", string(body)),
			zap.String("ip", c.ClientIP()),
			zap.Int("content_length", c.Writer.Size()),
		}

		var skip map[string]struct{}
		if length := len(options.skipPath); length > 0 {
			skip := make(map[string]struct{}, length)
			for _, path := range options.skipPath {
				skip[path] = struct{}{}
			}
		}

		// 当有skip path 或者 response body 过大时 log 中不记录 response
		if _, ok := skip[path]; ok || c.Writer.Size() >= 500 {
			fields = append(fields, zap.String("response", bodyWriter.body.String()))
		}

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			logger.Warn(msg, fields...)
		case c.Writer.Status() >= http.StatusInternalServerError:
			fields = append(fields, zap.String("errs", c.Errors.String()))
			logger.Error(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}
