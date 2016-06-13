package peppy

import (
	"database/sql"
	"strconv"

	"git.zxq.co/ripple/rippleapi/common"

	"github.com/gin-gonic/gin"
)

var defaultResponse = []struct{}{}

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
func genmodei(m string) int {
	v := common.Int(m)
	if v > 3 || v < 0 {
		v = 0
	}
	return v
}

func genUser(c *gin.Context, db *sql.DB) (string, string) {
	var whereClause string
	var p string

	// used in second case of switch
	s, err := strconv.Atoi(c.Query("u"))

	switch {
	// We know for sure that it's an username.
	case c.Query("type") == "string":
		whereClause = "users.username = ?"
		p = c.Query("u")
	// It could be an user ID, so we look for an user with that username first.
	case err == nil:
		err = db.QueryRow("SELECT id FROM users WHERE id = ? LIMIT 1", s).Scan(&p)
		if err == sql.ErrNoRows {
			// If no user with that userID were found, assume username.
			p = c.Query("u")
			whereClause = "users.username = ?"
		} else {
			// An user with that userID was found. Thus it's an userID.
			whereClause = "users.id = ?"
		}
	// u contains letters, so it's an username.
	default:
		p = c.Query("u")
		whereClause = "users.username = ?"
	}
	return whereClause, p
}
