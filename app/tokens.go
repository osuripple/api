package app

import (
	"database/sql"

	"github.com/osuripple/api/common"
)

// GetTokenFull retrieves an user ID and their token privileges knowing their API token.
func GetTokenFull(token string, db *sql.DB) (common.Token, bool) {
	var uid int
	var privs int
	err := db.QueryRow("SELECT user, privileges FROM tokens WHERE token = ? LIMIT 1", token).Scan(&uid, &privs)
	switch {
	case err == sql.ErrNoRows:
		return common.Token{}, false
	case err != nil:
		panic(err)
	default:
		return common.Token{
			Value:      token,
			UserID:     uid,
			Privileges: common.Privileges(privs),
		}, true
	}
}
