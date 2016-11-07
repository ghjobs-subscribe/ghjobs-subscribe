package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoas/stats"
)

// Stats provides response time, status code count, etc.
var Stats = stats.New()

func main() {
	r := gin.Default()

	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			beginning, recorder := Stats.Begin(c.Writer)
			c.Next()
			Stats.End(beginning, recorder)
		}
	}())

	r.LoadHTMLGlob("templates/*")

	r.POST("/subscribe", subscribeHandler)
	r.GET("/subscribe/verify", subscribeVerifyHandler)
	r.POST("/manage", manageHandler)
	r.POST("manage/update", manageUpdateHandler)
	r.POST("/unsubscribe", unsubscribeHandler)
	r.GET("/unsubscribe/verify", unsubscribeVerifyHandler)
	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, Stats.Data())
	})
	r.Run(":8080")
}
