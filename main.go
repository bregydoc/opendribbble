package main

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var ShotsDB *storm.DB
var CurrentShots []*GenericShot

func init() {

	var err error
	ShotsDB, err = storm.Open("shots.db")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(10 * time.Minute)

	CurrentShots, err = FetchAndUpdateShotsOnDB()
	if err != nil {
		panic(err)
	}

	go func() {
		for t := range ticker.C {
			log.Println(t.String())
			CurrentShots, err = FetchAndUpdateShotsOnDB()
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

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
