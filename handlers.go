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
		log.Panicf("error initializing DB: %v\n", err)
	}
	defer i.DB.Close()

	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")

	if _, err := mail.ParseAddress(email); len(email) == 0 || len(email) > 254 || err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "That email doesn't seem like a valid one.",
		})
	} else {
		ok, err := i.checkBucketExists(email)
		if err != nil {
			log.Panicf("error viewing bucket: %v\n", err)
		}
		if !ok {
			c.JSON(200, gin.H{
				"success": false,
				"message": "A subscription with this email already exists.",
			})
		} else {
			ok, err := i.createUserBucket(email)
			if err != nil {
				log.Panicf("error updating bucket: %v\n", err)
			}
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

func unsubscribeHandler(c *gin.Context) {
	i := impl{}
	err := i.initDB()
	if err != nil {
		log.Panicf("error initializing DB: %v\n", err)
	}
	defer i.DB.Close()

	email := c.PostForm("email")
	c.Header("Access-Control-Allow-Origin", "*")

	ok, err := i.checkBucketExists(email)
	if err != nil {
		log.Panicf("error viewing bucket: %v\n", err)
	}
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
