package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/subscribe", subscribeHandler)
	// r.POST("/manage", manageHander)
	r.POST("/unsubscribe", unsubscribeHandler)
	r.Run(":8080")
}
