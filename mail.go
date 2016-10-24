package main

import (
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ghjsCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func sendActivationMail(to string) error {
	from := "mail@ghjobssubscribe.com"
	password := os.Getenv("GHJS_NOREPLY_PASSWORD")

	body, err := generateActivationBody(to)
	if err != nil {
		return err
	}

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		"Content-type: text/html" + "\r\n" +
		"Subject: GHJobs Subscribe: Please Confrim Subscription" + "\r\n\r\n" +
		body + "\r\n"

	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func generateActivationBody(email string) (string, error) {
	claims := &ghjsCustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			Issuer:    "ghjobssubscribe",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("GHJS_SECRET_STRING")))

	body := "<p>Click on the following link to activate your account: <br>" + getActivationURL(tokenString) + "<br>The link will expire in 30 minutes.</p>"

	return body, err
}

func getActivationURL(tokenString string) string {
	return "https://api.ghjobssubscribe.com/subscribe/verify?token=" + tokenString
}
