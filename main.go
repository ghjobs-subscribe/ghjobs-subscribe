package main

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/subscribe", subscribeHandler)

	r.Run(":8080")
}

func subscribeHandler(c *gin.Context) {
	email := c.PostForm("email")
	if len(email) == 0 {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Looks like you forgot to enter your email.",
		})
	} else if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
		c.JSON(200, gin.H{
			"success": false,
			"message": "That email doesn't seem like a valid one.",
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Subscribed! Check your email for confirmation.",
		})
	}
}
