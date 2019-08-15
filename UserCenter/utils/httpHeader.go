package utils

import "github.com/gin-gonic/gin"

// set photo content type header
func SetContentTypeImage(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "image/png;")
}