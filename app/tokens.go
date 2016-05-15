package app

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	"git.zxq.co/ripple/rippleapi/common"
)

// GetTokenFull retrieves an user ID and their token privileges knowing their API token.
func GetTokenFull(token string, db *sql.DB) (common.Token, bool) {
	var uid int
	var privs int
	var priv8 bool
	err := db.QueryRow("SELECT user, privileges, private FROM tokens WHERE token = ? LIMIT 1", fmt.Sprintf("%x", md5.Sum([]byte(token)))).Scan(&uid, &privs, &priv8)
	if priv8 {
		privs = common.PrivilegeRead | common.PrivilegeReadConfidential | common.PrivilegeWrite
	}
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
