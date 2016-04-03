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
