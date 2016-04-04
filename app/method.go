package app

import (
	"database/sql"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/osuripple/api/common"
)

// Method wraps an API method to a HandlerFunc.
func Method(f func(md common.MethodData) common.Response, db *sql.DB, privilegesNeeded ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
		}
		c.Request.Body.Close()

		token := ""
		switch {
		case c.Request.Header.Get("X-Ripple-Token") != "":
			token = c.Request.Header.Get("X-Ripple-Token")
		case c.Query("token") != "":
			token = c.Query("token")
		case c.Query("k") != "":
			token = c.Query("k")
		}

		md := common.MethodData{
			DB:          db,
			RequestData: data,
			C:           c,
		}
		if token != "" {
			tokenReal, exists := GetTokenFull(token, db)
			if exists {
				md.User = tokenReal
			}
		}

		missingPrivileges := 0
		for _, privilege := range privilegesNeeded {
			if int(md.User.Privileges)&privilege == 0 {
				missingPrivileges |= privilege
			}
		}
		if missingPrivileges != 0 {
			c.IndentedJSON(401, common.Response{
				Code:    401,
				Message: "You don't have the privilege(s): " + common.Privileges(missingPrivileges).String() + ".",
			})
			return
		}

		resp := f(md)
		if resp.Code == 0 {
			resp.Code = 500
		}
		if _, exists := c.GetQuery("pls200"); exists {
			c.IndentedJSON(200, resp)
		} else {
			c.IndentedJSON(resp.Code, resp)
		}
	}
}
