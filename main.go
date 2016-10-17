package main

import (
	"net/mail"

	"github.com/gin-gonic/gin"
)

func main() {
	i := impl{}
	i.initDB()
	defer i.DB.Close()

	r := gin.Default()
	r.POST("/subscribe", i.subscribeHandler)
	r.POST("/unsubscribe", i.unsubscribeHandler)
	r.Run(":8080")
}

func (i *impl) subscribeHandler(c *gin.Context) {
	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")
	if _, err := mail.ParseAddress(email); len(email) == 0 || len(email) > 254 || err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "That email doesn't seem like a valid one.",
		})
	} else {
		ok := i.checkBucketExists(email)
		if !ok {
			c.JSON(200, gin.H{
				"success": false,
				"message": "A subscription with this email already exists.",
			})
		} else {
			ok := i.createUserBucket(email)
			if !ok {
				c.JSON(200, gin.H{
					"success": false,
					"message": "An internal error occured. Please try again later.",
				})
			} else {
				c.JSON(200, gin.H{
					"success": true,
					"message": "All set! Check your email for subscription confirmation.",
				})
			}
		}
	}
}

func (i *impl) unsubscribeHandler(c *gin.Context) {
	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")
	ok := i.checkBucketExists(email)
	if !ok {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Sad that you are leaving. Check your email for unsubscription confirmation.",
		})

	} else {
		c.JSON(200, gin.H{
			"success": false,
			"message": "A subscription with this email does not exist.",
		})
	}
}
