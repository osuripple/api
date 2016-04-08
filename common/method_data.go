package common

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// MethodData is a struct containing the data passed over to an API method.
type MethodData struct {
	User        Token
	DB          *sql.DB
	RequestData []byte
	C           *gin.Context
}

// Err logs an error into gin.
func (md MethodData) Err(err error) {
	md.C.Error(err)
}

// ID retrieves the Token's owner user ID.
func (md MethodData) ID() int {
	return md.User.UserID
}
