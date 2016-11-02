package main

import (
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GHJSCustomClaims holds claims for generating JWT tokens
type GHJSCustomClaims struct {
	Email     string `json:"email"`
	Subscribe string `json:"subscribe"`
	jwt.StandardClaims
}

func sendVerificationMail(to string, action bool) error {
	from := "mail@ghjobssubscribe.com"
	password := os.Getenv("GHJS_NOREPLY_PASSWORD")

	act := "Subscription"
	if action != true {
		act = "Unsubscription"
	}

	body, err := generateActivationBody(to, action)
	if err != nil {
		return err
	}

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		"Content-type: text/html" + "\r\n" +
		"Subject: GHJobs Subscribe: Please Confrim " + act + "\r\n\r\n" +
		body + "\r\n"

	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func generateActivationBody(email string, action bool) (string, error) {
	subscribe := "true"
	if !action {
		subscribe = "false"
	}

	claims := GHJSCustomClaims{
		email,
		subscribe,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "ghjobssubscribe",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("GHJS_SECRET_STRING")))

	body := "<p>Click on the following link to activate your account: <br>" + getActivationURL(tokenString) + "<br>The link will expire in 24 hours.</p>"
	if !action {
		body = "<p>Click on the following link to deactivate your account: <br>" + getDeactivationURL(tokenString) + "<br>The link will expire in 24 hours.</p>"
	}

	return body, err
}

func getActivationURL(tokenString string) string {
	return "https://api.ghjobssubscribe.com/subscribe/verify?token=" + tokenString
}

func getDeactivationURL(tokenString string) string {
	return "https://api.ghjobssubscribe.com/unsubscribe/verify?token=" + tokenString
}

// func getActivationURL(tokenString string) string {
// 	return "http://localhost:8080/subscribe/verify?token=" + tokenString
// }

// func getDeactivationURL(tokenString string) string {
// 	return "http://localhost:8080/unsubscribe/verify?token=" + tokenString
// }
