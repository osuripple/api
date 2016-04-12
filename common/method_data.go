package common

import (
	"database/sql"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// MethodData is a struct containing the data passed over to an API method.
type MethodData struct {
	User        Token
	DB          *sql.DB
	RequestData RequestData
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

// RequestData is the body of a request. It is wrapped into this type
// to implement the Unmarshal function, which is just a shorthand to
// json.Unmarshal.
type RequestData []byte

// Unmarshal json-decodes Requestdata into a value. Basically a
// shorthand to json.Unmarshal.
func (r RequestData) Unmarshal(into interface{}) error {
	return json.Unmarshal([]byte(r), into)
}
