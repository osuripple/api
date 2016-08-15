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
	var privs uint64
	var priv8 bool
	err := db.QueryRow("SELECT id, user, privileges, private FROM tokens WHERE token = ? LIMIT 1",
		fmt.Sprintf("%x", md5.Sum([]byte(token)))).
		Scan(
			&t.ID, &t.UserID, &privs, &priv8,
		)
	if priv8 {
		privs = common.PrivilegeReadConfidential | common.PrivilegeWrite
	}
	t.Privileges = common.Privileges(privs)
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
