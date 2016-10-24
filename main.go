package main

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	templatePath := path.Join("$GOPATH/src/github.com/ghjobs-subscribe/ghjobs-subscribe/templates")
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tmpl", templatePath))
	r.LoadHTMLFiles(files...)

	r.POST("/subscribe", subscribeHandler)
	r.GET("/subscribe/verify", subscribeVerifyHandler)
	// r.POST("/manage", manageHander)
	r.POST("/unsubscribe", unsubscribeHandler)
	r.Run(":8080")
}
