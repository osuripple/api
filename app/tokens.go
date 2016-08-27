package app

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"git.zxq.co/ripple/rippleapi/common"
)

// GetTokenFull retrieves an user ID and their token privileges knowing their API token.
func GetTokenFull(token string, db *sqlx.DB) (common.Token, bool) {
	var t common.Token
	var (
		tokenPrivsRaw uint64
		userPrivsRaw  uint64
	)
	var priv8 bool
	err := db.QueryRow(`SELECT 
	t.id, t.user, t.privileges, t.private, u.privileges
FROM tokens t
LEFT JOIN users u ON u.id = t.user
WHERE token = ? LIMIT 1`,
		fmt.Sprintf("%x", md5.Sum([]byte(token)))).
		Scan(
			&t.ID, &t.UserID, &tokenPrivsRaw, &priv8, &userPrivsRaw,
		)
	if priv8 {
		tokenPrivsRaw = common.PrivilegeReadConfidential | common.PrivilegeWrite
	}
	t.TokenPrivileges = common.Privileges(tokenPrivsRaw)
	t.UserPrivileges = common.UserPrivileges(userPrivsRaw)
	switch {
	case err == sql.ErrNoRows:
		return common.Token{}, false
	case err != nil:
		panic(err)
	default:
		t.Value = token
		return t, true
	}
}
