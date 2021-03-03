/**
 * @Time: 2021/2/28 8:13 下午
 * @Author: varluffy
 * @Description: v1.router
 */

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/varluffy/ginx/example/api/middleware"
	"github.com/varluffy/ginx/example/internal/service"
)

func RegisterRouter(router *gin.Engine, service *service.Service, token middleware.Token) {
	a := router.Group("auth")
	a.POST("login", service.Auth.Login)
	user := router.Group("user")
	authMiddleware := middleware.NewAuth(token).Auth()
	user.Use(authMiddleware)
	user.GET("info", func(c *gin.Context) {
		c.JSON(200, gin.H{"userId": middleware.FromUserId(c)})
	})
}
