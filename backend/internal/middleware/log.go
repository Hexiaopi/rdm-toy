package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type ResponseWithRecorder struct {
	gin.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *ResponseWithRecorder) WriteHeader(statusCode int) {
	rec.ResponseWriter.WriteHeader(statusCode)
	rec.statusCode = statusCode
}

func (rec *ResponseWithRecorder) Write(d []byte) (n int, err error) {
	n, err = rec.ResponseWriter.Write(d)
	if err != nil {
		return
	}
	rec.body.Write(d)
	return
}

// Logger 日志记录
func Logger(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		start := time.Now()
		//记录请求包
		buf, _ := io.ReadAll(c.Request.Body)
		rdr := io.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = rdr //rewrite

		log.WithFields(log.Fields{
			XRequestIDKey: c.GetString(XRequestIDKey),
			"path":        c.Request.URL.Path,
			"param":       c.Request.URL.RawQuery,
			"method":      c.Request.Method,
			"reqpkg":      string(buf),
		}).Info("receive request")

		//记录返回包
		wc := &ResponseWithRecorder{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
			body:           bytes.Buffer{},
		}
		c.Writer = wc

		c.Next()

		defer func() { //日志记录扫尾工作
			entry := log.WithFields(log.Fields{
				XRequestIDKey: c.GetString(XRequestIDKey),
				"path":        c.Request.URL.Path,
				"status":      wc.statusCode,
				"respkg":      wc.body.String(),
				"utm":         time.Since(start).String(),
			})
			if errors := c.Errors.Errors(); len(errors) > 0 {
				entry.Error(errors)
			} else {
				entry.Info("done request")
			}
		}()
	}
}
