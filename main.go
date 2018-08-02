package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gincors := cors.DefaultConfig()
	gincors.AllowAllOrigins = true
	r.Use(cors.New(gincors))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"about":  "this is a unofficial api for get popular shots from dribbble.",
			"source": "https://github.com/bregydoc/opendribbble",
			"I":      "❤️ open source",
		})
	})
	LinkAPI("/api", r)
	r.Run(":4700")
}
