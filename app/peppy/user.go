// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/thehowl/go-osuapi"
)

// GetUser retrieves general user information.
func GetUser(c *gin.Context, db *sql.DB) {
	if c.Query("u") == "" {
		c.JSON(200, []struct{}{})
		return
	}
	var user osuapi.User
	whereClause, p := genUser(c, db)

	mode := genmode(c.Query("m"))

	fmt.Println(whereClause, p)

	var display bool
	err := db.QueryRow(fmt.Sprintf(
		`SELECT
			users.id, users.username,
			users_stats.playcount_%s, users_stats.ranked_score_%s, users_stats.total_score_%s,
			leaderboard_%s.position, users_stats.pp_%s, users_stats.avg_accuracy_%s,
			users_stats.country, users_stats.show_country
		FROM users
		LEFT JOIN users_stats ON users_stats.id = users.id
		LEFT JOIN leaderboard_%s ON leaderboard_%s.user = users.id
		%s
		LIMIT 1`,
		mode, mode, mode, mode, mode, mode, mode, mode, whereClause,
	), p).Scan(
		&user.UserID, &user.Username,
		&user.Playcount, &user.RankedScore, &user.TotalScore,
		&user.Rank, &user.PP, &user.Accuracy,
		&user.Country, &display,
	)
	if err != nil {
		c.JSON(200, []struct{}{})
		return
	}
	if !display {
		user.Country = "XX"
	}

	c.JSON(200, []osuapi.User{user})
}
