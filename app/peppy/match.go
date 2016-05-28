// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// GetMatch retrieves general match information.
func GetMatch(c *gin.Context, db *sql.DB) {
	c.JSON(200, defaultResponse)
}
