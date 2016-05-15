package peppy

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func genmode(m string) string {
	switch m {
	case "1":
		m = "taiko"
	case "2":
		m = "ctb"
	case "3":
		m = "mania"
	default:
		m = "std"
	}
	return m
}

func genUser(c *gin.Context, db *sql.DB) (string, string) {
	var whereClause string
	var p string

	// used in second case of switch
	_, err := strconv.Atoi(c.Query("u"))

	switch {
	// We know for sure that it's an username.
	case c.Query("type") == "string":
		whereClause = "WHERE users.username = ?"
		p = c.Query("u")
	// It could be an user ID, so we look for an user with that username first.
	case err == nil:
		err = db.QueryRow("SELECT id FROM users WHERE username = ? LIMIT 1", c.Query("u")).Scan(&p)
		// If there is an error, that means u is an userID.
		// If there is none, p will automatically have become the user id retrieved from the database
		// in the last query.
		if err == sql.ErrNoRows {
			p = c.Query("u")
		}
		whereClause = "WHERE users.id = ?"
	// u contains letters, so it's an username.
	default:
		p = c.Query("u")
		whereClause = "WHERE users.username = ?"
	}
	return whereClause, p
}
