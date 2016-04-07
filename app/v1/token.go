package v1

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/osuripple/api/common"
	"golang.org/x/crypto/bcrypt"
)

type tokenNewInData struct {
	// either username or userid must be given in the request.
	// if none is given, the request is trashed.
	Username    string `json:"username"`
	UserID      int    `json:"id"`
	Password    string `json:"password"`
	Privileges  int    `json:"privileges"`
	Description string `json:"description"`
}

type tokenNewOutData struct {
	Username   string `json:"username"`
	ID         int    `json:"id"`
	Privileges int    `json:"privileges"`
	Token      string `json:"token,omitempty"`
	Banned     bool   `json:"banned"`
}

// TokenNewPOST is the handler for POST /token/new.
func TokenNewPOST(md common.MethodData) (r common.Response) {
	data := tokenNewInData{}
	err := json.Unmarshal(md.RequestData, &data)
	if err != nil {
		r = ErrBadJSON
		return
	}

	var miss []string
	if data.Username == "" && data.UserID == 0 {
		miss = append(miss, "username|id")
	}
	if data.Password == "" {
		miss = append(miss, "password")
	}
	if len(miss) != 0 {
		r = ErrMissingField(miss...)
		return
	}

	var q *sql.Row
	const base = "SELECT id, username, rank, password_md5, password_version, allowed FROM users "
	if data.UserID != 0 {
		q = md.DB.QueryRow(base+"WHERE id = ? LIMIT 1", data.UserID)
	} else {
		q = md.DB.QueryRow(base+"WHERE username = ? LIMIT 1", data.Username)
	}

	ret := tokenNewOutData{}
	var (
		rank      int
		pw        string
		pwVersion int
		allowed   int
	)

	err = q.Scan(&ret.ID, &ret.Username, &rank, &pw, &pwVersion, &allowed)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "No user with that username/id was found."
		return
	case err != nil:
		md.C.Error(err)
		r = Err500
		return
	}

	if pwVersion == 1 {
		r.Code = 418 // Teapots!
		r.Message = "That user still has a password in version 1. Unfortunately, in order for the API to check for the password to be OK, the user has to first log in through the website."
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(fmt.Sprintf("%x", md5.Sum([]byte(data.Password))))); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			r.Code = 403
			r.Message = "That password doesn't match!"
			return
		}
		md.C.Error(err)
		r = Err500
		return
	}
	if allowed == 0 {
		r.Code = 200
		r.Message = "That user is banned."
		ret.Banned = true
		r.Data = ret
		return
	}
	ret.Privileges = int(common.Privileges(data.Privileges).CanOnly(rank))

	var (
		tokenStr string
		tokenMD5 string
	)
	for {
		tokenStr = common.RandomString(32)
		tokenMD5 = fmt.Sprintf("%x", md5.Sum([]byte(tokenStr)))
		ret.Token = tokenStr
		id := 0

		err := md.DB.QueryRow("SELECT id FROM tokens WHERE token=? LIMIT 1", tokenMD5).Scan(&id)
		if err == sql.ErrNoRows {
			break
		}
		if err != nil {
			md.C.Error(err)
			r = Err500
			return
		}
	}
	_, err = md.DB.Exec("INSERT INTO tokens(user, privileges, description, token) VALUES (?, ?, ?, ?)", ret.ID, ret.Privileges, data.Description, tokenMD5)
	if err != nil {
		md.C.Error(err)
		r = Err500
		return
	}

	r.Code = 200
	r.Data = ret
	return
}
