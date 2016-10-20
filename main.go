package main

import (
	"log"
	"os"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/subscribe", subscribeHandler)
	// r.POST("/manage", manageHander)
	r.POST("/unsubscribe", unsubscribeHandler)

	err := endless.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
	log.Println("ghjobs-subscribe server stopped on 8080")

	os.Exit(0)
}
