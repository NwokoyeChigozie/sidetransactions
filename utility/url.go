package utility

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHeader(c *gin.Context, key string) string {
	header := ""
	if c.GetHeader(key) != "" {
		header = c.GetHeader(key)
	} else if c.GetHeader(strings.ToLower(key)) != "" {
		header = c.GetHeader(strings.ToLower(key))
	} else if c.GetHeader(strings.ToUpper(key)) != "" {
		header = c.GetHeader(strings.ToUpper(key))
	} else if c.GetHeader(strings.Title(key)) != "" {
		header = c.GetHeader(strings.Title(key))
	}
	return header
}
