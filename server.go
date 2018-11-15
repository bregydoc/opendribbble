package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func Shuffle(shots []*GenericShot) []*GenericShot {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]*GenericShot, len(shots))
	perm := r.Perm(len(shots))
	for i, randIndex := range perm {
		ret[i] = shots[randIndex]
	}
	return ret
}

// LinkAPI ...
func LinkAPI(base string, r *gin.Engine) {
	api := r.Group(base)
	api.GET("/popular_shots", func(c *gin.Context) {
		shots := CurrentShots
		shots = Shuffle(shots)
		c.JSON(http.StatusOK, shots)
		//feed, _ := GetFeedFromKeyword("machine learning", map[string]string{
		//	"max_results": "30",
		//})
		//c.JSON(http.StatusOK, feed.Papers)
	})
}
