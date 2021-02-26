/**
 * @Time: 2021/2/24 6:25 下午
 * @Author: varluffy
 * @Description: server test
 */

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/log"
	router2 "github.com/varluffy/ginx/transport/http/router"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"
)

type testData struct {
	Path string `json:"path"`
}

func TestServer(t *testing.T) {
	fn := gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(200, &testData{Path: c.Request.RequestURI})
		return
	})
	logger := log.NewLogger(log.WithConsoleEncoder())
	router := router2.NewRouter(router2.WithLogger(logger))
	group := router.Group("/test")
	{
		group.GET("/", fn)
		group.HEAD("/index?a=1&b=2", fn)
		group.OPTIONS("/home", fn)
		group.PUT("/products/:id", fn)
		group.POST("/products/:id", fn)
		group.PATCH("/products/:id", fn)
		group.DELETE("/products/:id", fn)
	}
	srv := NewServer(Router(router), Logger(logger))
	time.AfterFunc(time.Second, func() {
		defer srv.Stop()
		testClient(t, srv)
	})
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
}

func testClient(t *testing.T, srv *Server) {
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/test/"},
		{"PUT", "/test/products/1?a=1&b=2"},
		{"POST", "/test/products/2"},
		{"PATCH", "/test/products/3"},
		{"DELETE", "/test/products/4"},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, req)
		result := w.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)
		t.Logf("%s", body)
	}
}
