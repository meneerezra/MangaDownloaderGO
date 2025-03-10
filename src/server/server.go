package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartServer(router *gin.Engine) {
	router.LoadHTMLGlob("server/templates/*.html")
	router.Static("/static", "server/templates/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	err := router.Run(":80")
	if err != nil {
		return 
	}
}
