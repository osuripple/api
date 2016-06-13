package app

import (
	"encoding/json"
	"io/ioutil"
	"regexp"

	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
)

// Method wraps an API method to a HandlerFunc.
func Method(f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		initialCaretaker(c, f, privilegesNeeded...)
	}
}

func initialCaretaker(c *gin.Context, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
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
	default:
		token, _ = c.Cookie("X-Ripple-Token")
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
		c.IndentedJSON(401, common.SimpleResponse(401, "You don't have the privilege(s): "+common.Privileges(missingPrivileges).String()+"."))
		return
	}

	resp := f(md)
	if resp.GetCode() == 0 {
		// Dirty hack to set the code
		type setCoder interface {
			SetCode(int)
		}
		if newver, can := resp.(setCoder); can {
			newver.SetCode(500)
		}
	}

	if _, exists := c.GetQuery("pls200"); exists {
		c.Writer.WriteHeader(200)
	} else {
		c.Writer.WriteHeader(resp.GetCode())
	}

	if _, exists := c.GetQuery("callback"); exists {
		c.Header("Content-Type", "application/javascript; charset=utf-8")
	} else {
		c.Header("Content-Type", "application/json; charset=utf-8")
	}

	mkjson(c, resp)
}

// Very restrictive, but this way it shouldn't completely fuck up.
var callbackJSONP = regexp.MustCompile(`^[a-zA-Z_\$][a-zA-Z0-9_\$]*$`)

// mkjson auto indents json, and wraps json into a jsonp callback if specified by the request.
// then writes to the gin.Context the data.
func mkjson(c *gin.Context, data interface{}) {
	exported, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		c.Error(err)
		exported = []byte(`{ "code": 500, "message": "something has gone really really really really really really wrong.", "data": null }`)
	}
	cb := c.Query("callback")
	willcb := cb != "" &&
		len(cb) < 100 &&
		callbackJSONP.MatchString(cb)
	if willcb {
		c.Writer.Write([]byte("/**/ typeof " + cb + " === 'function' && " + cb + "("))
	}
	c.Writer.Write(exported)
	if willcb {
		c.Writer.Write([]byte(");"))
	}
}
