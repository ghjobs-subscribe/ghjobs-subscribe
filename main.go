package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.POST("/subscribe", subscribeHandler)
	r.GET("/subscribe/verify", subscribeVerifyHandler)
	// r.POST("/manage", manageHander)
	r.POST("/unsubscribe", unsubscribeHandler)
	r.Run(":8080")
}
