package app

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware for gin that takes care of calls to c.Error().
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors.Errors()
		if len(errs) != 0 {
			color.Red("!!! ERRORS OCCURRED !!!")
			color.Red("==> %s %s", c.Request.Method, c.Request.URL.Path)
			for _, err := range errs {
				color.Red("===> %s", err)
			}
		}
	}
}
