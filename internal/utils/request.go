package utils

import "github.com/gin-gonic/gin"

func BuildBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + c.Request.URL.Path
}
