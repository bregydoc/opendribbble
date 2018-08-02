package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LinkAPI ...
func LinkAPI(base string, r *gin.Engine) {
	api := r.Group(base)
	api.GET("/popular_shots", func(c *gin.Context) {
		shots := CollectAllPopularShots()
		c.JSON(http.StatusOK, shots)
	})
}
