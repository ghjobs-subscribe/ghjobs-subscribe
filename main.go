package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.POST("/subscribe", subscribeHandler)
	r.GET("/subscribe/verify", subscribeVerifyHandler)
	r.POST("/manage", manageHandler)
	r.POST("manage/update", manageUpdateHandler)
	r.POST("/unsubscribe", unsubscribeHandler)
	r.GET("/unsubscribe/verify", unsubscribeVerifyHandler)
	r.Run(":8080")
}
