package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
				err := sendVerificationMail(email, true)
				if err != nil {
					log.Printf("error sending email: %v\n", err)
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
}

func subscribeVerifyHandler(c *gin.Context) {
	i := impl{}
	err := i.initDB()
	if err != nil {
		log.Printf("error initializing DB: %v\n", err)
		c.HTML(200, "verify.tmpl", gin.H{
			"success": false,
			"message": "An internal error occured. Please try again later.",
		})
	}
	defer i.DB.Close()

	token, err := jwt.ParseWithClaims(c.Query("token"), &GHJSCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("GHJS_SECRET_STRING")), nil
	})
	if err != nil {
		c.HTML(200, "verify.tmpl", gin.H{
			"success":        false,
			"showmanagelink": false,
			"message":        "Invalid signing algorithm.",
		})
	} else {
		if claims, ok := token.Claims.(*GHJSCustomClaims); ok && token.Valid {
			if (claims.StandardClaims.Issuer == "ghjobssubscribe") && (claims.Subscribe == "true") && (claims.StandardClaims.ExpiresAt > time.Now().Unix()) {
				if !i.checkUserSubscription(claims.Email) {
					err := i.changeUserSubscription(claims.Email, "true")
					if err != nil {
						c.HTML(200, "verify.tmpl", gin.H{
							"success":        false,
							"email":          claims.Email,
							"showmanagelink": false,
							"message":        "An internal error occured. Click on the link (sent to your email) again.",
						})
					} else {
						c.HTML(200, "verify.tmpl", gin.H{
							"success":        true,
							"email":          claims.Email,
							"showmanagelink": true,
							"message":        "Your subscription is now active.",
						})
					}
				} else {
					c.HTML(200, "verify.tmpl", gin.H{
						"success":        false,
						"email":          claims.Email,
						"showmanagelink": true,
						"message":        "Your subscription is already active.",
					})
				}
			} else {
				c.HTML(200, "verify.tmpl", gin.H{
					"success":        false,
					"email":          claims.Email,
					"showmanagelink": true,
					"message":        "Invalid token.",
				})
			}
		} else {
			c.HTML(200, "verify.tmpl", gin.H{
				"success":        false,
				"showmanagelink": true,
				"message":        "Invalid token.",
			})
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
			err := sendVerificationMail(email, false)
			if err != nil {
				log.Printf("error sending email: %v\n", err)
				c.JSON(200, gin.H{
					"success": false,
					"message": "An internal error occured. Please try again later.",
				})
			} else {
				c.JSON(200, gin.H{
					"success": true,
					"message": "Sad that you are leaving. Check your email for a confirmation.",
				})
			}
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"message": "Your subscription is already inactive.",
			})
		}
	} else {
		c.JSON(200, gin.H{
			"success": false,
			"message": "A subscription with this email does not exist.",
		})
	}
}

func unsubscribeVerifyHandler(c *gin.Context) {
	i := impl{}
	err := i.initDB()
	if err != nil {
		log.Printf("error initializing DB: %v\n", err)
		c.HTML(200, "verify.tmpl", gin.H{
			"success": false,
			"message": "An internal error occured. Please try again later.",
		})
	}
	defer i.DB.Close()

	token, err := jwt.ParseWithClaims(c.Query("token"), &GHJSCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("GHJS_SECRET_STRING")), nil
	})
	if err != nil {
		c.HTML(200, "verify.tmpl", gin.H{
			"success":        false,
			"showmanagelink": false,
			"message":        "Invalid signing algorithm.",
		})
	} else {
		if claims, ok := token.Claims.(*GHJSCustomClaims); ok && token.Valid {
			if (claims.StandardClaims.Issuer == "ghjobssubscribe") && (claims.Subscribe == "false") && (claims.StandardClaims.ExpiresAt > time.Now().Unix()) {
				if i.checkUserSubscription(claims.Email) {
					err := i.changeUserSubscription(claims.Email, "false")
					if err != nil {
						c.HTML(200, "verify.tmpl", gin.H{
							"success":        false,
							"email":          claims.Email,
							"showmanagelink": false,
							"message":        "An internal error occured. Click on the link (sent to your email) again.",
						})
					} else {
						c.HTML(200, "verify.tmpl", gin.H{
							"success":        true,
							"email":          claims.Email,
							"showmanagelink": true,
							"message":        "Your subscription is now inactive.",
						})
					}
				} else {
					c.HTML(200, "verify.tmpl", gin.H{
						"success":        false,
						"email":          claims.Email,
						"showmanagelink": true,
						"message":        "Your subscription is already inactve.",
					})
				}
			} else {
				c.HTML(200, "verify.tmpl", gin.H{
					"success":        false,
					"email":          claims.Email,
					"showmanagelink": true,
					"message":        "Invalid token.",
				})
			}
		} else {
			c.HTML(200, "verify.tmpl", gin.H{
				"success":        false,
				"showmanagelink": true,
				"message":        "Invalid token.",
			})
		}
	}
}
