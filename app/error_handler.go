package app

import (
	"fmt"

	"git.zxq.co/ripple/schiavolib"
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
			var out string
			out += fmt.Sprintf("==> %s %s\n", c.Request.Method, c.Request.URL.Path)
			for _, err := range errs {
				out += fmt.Sprintf("===> %s\n", err)
			}
			color.Red(out)
			go schiavo.Bunker.Send("Errors occurred:\n```\n" + out + "```")
		}
	}
}
