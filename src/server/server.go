package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartServer(router *gin.Engine) {
	router.LoadHTMLGlob("../templates/*.html")
	router.Static("/static", "../templates/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run(":80")
}
