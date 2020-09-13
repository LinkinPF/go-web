package routes

import (
	"github.com/gin-gonic/gin"
	"go-web/10-arch/logger"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})
	return r
}
