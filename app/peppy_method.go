package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// PeppyMethod generates a method for the peppyapi
func PeppyMethod(a func(c *gin.Context, db *sql.DB)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// I have no idea how, but I manged to accidentally string the first 4
		// letters of the alphabet into a single function call.
		a(c, db)
	}
}
