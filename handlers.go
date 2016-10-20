package main

import (
	"log"
	"net/mail"

	"github.com/gin-gonic/gin"
)

func subscribeHandler(c *gin.Context) {
	i := impl{}
	err := i.initDB()
	if err != nil {
		log.Printf("error initializing DB: %v\n", err)
		c.JSON(200, gin.H{
			"success": false,
			"message": "An internal error occured. Please try again later.",
		})
	}
	defer i.DB.Close()

	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")

	_, err = mail.ParseAddress(email)
	if len(email) == 0 || len(email) > 254 || err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "That email doesn't seem like a valid one.",
		})
	} else {
		ok := i.checkUserExists(email)
		if ok != false {
			c.JSON(200, gin.H{
				"success": false,
				"message": "A subscription with this email already exists.",
			})
		} else {
			err := i.createUserProfile(email)
			if err != nil {
				log.Printf("error updating bucket: %v\n", err)
				c.JSON(200, gin.H{
					"success": false,
					"message": "An internal error occured. Please try again later.",
				})
			} else {
				c.JSON(200, gin.H{
					"success": true,
					"message": "All set! Check your email for a confirmation.",
				})
			}
		}
	}
}

func unsubscribeHandler(c *gin.Context) {
	i := impl{}
	err := i.initDB()
	if err != nil {
		log.Panicf("error initializing DB: %v\n", err)
	}
	defer i.DB.Close()

	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")

	ok := i.checkUserExists(email)
	if ok != false {
		ok = i.checkUserSubscription(email)
		if ok != false {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Sad that you are leaving. Check your email for a confirmation.<br>If you have a minute, please send a message about what made you unsubscribe. Your feedback will be appreciated.",
			})
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"message": "Looks like you have not activated account yet.",
			})
		}
	} else {
		c.JSON(200, gin.H{
			"success": false,
			"message": "A subscription with this email does not exist.",
		})
	}
}
