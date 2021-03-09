/**
 * @Time: 2021/2/28 8:13 下午
 * @Author: varluffy
 */

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/varluffy/rich/example/api/middleware"
	"github.com/varluffy/rich/example/internal/service"
	"github.com/varluffy/rich/transport/http"
	"github.com/varluffy/rich/transport/http/gin/ginx"
)

var RouterSet = wire.NewSet(NewRouter)

type Router struct {
	middleware *middleware.Middleware
	service    *service.Service
	server     *http.Server
}

func NewRouter(server *http.Server, middleware *middleware.Middleware, service *service.Service) *Router {
	return &Router{
		middleware: middleware,
		service:    service,
		server:     server,
	}
}

func (r *Router) Register() {
	router := r.server.Router()
	auth := router.Group("auth")
	auth.POST("login", r.service.Auth.Login)
	user := router.Group("user").Use(r.middleware.Auth())
	user.POST("info", func(c *gin.Context) {
		ginx.Response(c, gin.H{"userId": middleware.FromUserId(c)})
	})
}
