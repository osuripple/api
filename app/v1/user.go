// Package v1 implements the first version of the Ripple API.
package v1

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/osuripple/api/common"
)

type userData struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	UsernameAKA    string    `json:"username_aka"`
	RegisteredOn   time.Time `json:"registered_on"`
	Rank           int       `json:"rank"`
	LatestActivity time.Time `json:"latest_activity"`
	Country        string    `json:"country"`
	Badges         []int     `json:"badges"`
}

// UserByIDGET is the API handler for GET /users/id/:id
func UserByIDGET(md common.MethodData) (r common.Response) {
	var err error
	var uid int
	uidStr := md.C.Param("id")
	if uidStr == "self" {
		uid = md.User.UserID
	} else {
		uid, err = strconv.Atoi(uidStr)
		if err != nil {
			r.Code = 400
			r.Message = fmt.Sprintf("%s ain't a number", uidStr)
			return
		}
	}

	query := `
SELECT users.id, users.username, register_datetime, rank,
	latest_activity, users_stats.username_aka, users_stats.badges_shown,
	users_stats.country, users_stats.show_country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.id=?
LIMIT 1`
	r = userPuts(md, md.DB.QueryRow(query, uid))
	return
}

// UserByNameGET is the API handler for GET /users/name/:name
func UserByNameGET(md common.MethodData) (r common.Response) {
	username := md.C.Param("name")

	query := `
SELECT users.id, users.username, register_datetime, rank,
	latest_activity, users_stats.username_aka, users_stats.badges_shown,
	users_stats.country, users_stats.show_country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.username=?
LIMIT 1`
	r = userPuts(md, md.DB.QueryRow(query, username))
	return
}

func userPuts(md common.MethodData, row *sql.Row) (r common.Response) {
	var err error
	var user userData

	registeredOn := int64(0)
	latestActivity := int64(0)
	var badges string
	var showcountry bool
	err = row.Scan(&user.ID, &user.Username, &registeredOn, &user.Rank, &latestActivity, &user.UsernameAKA, &badges, &user.Country, &showcountry)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "No such user was found!"
		return
	case err != nil:
		md.C.Error(err)
		r = Err500
		return
	}

	user.RegisteredOn = time.Unix(registeredOn, 0)
	user.LatestActivity = time.Unix(latestActivity, 0)

	badgesSl := strings.Split(badges, ",")
	for _, badge := range badgesSl {
		if badge != "" && badge != "0" {
			// We are ignoring errors because who really gives a shit if something's gone wrong on our end in this
			// particular thing, we can just silently ignore this.
			nb, err := strconv.Atoi(badge)
			if err == nil && nb != 0 {
				user.Badges = append(user.Badges, nb)
			}
		}
	}

	// If the user wants to stay anonymous, don't show their country.
	// This can be overriden if we have the ReadConfidential privilege and the user we are accessing is the token owner.
	if !(showcountry || (md.User.Privileges.HasPrivilegeReadConfidential() && user.ID == md.User.UserID)) {
		user.Country = "XX"
	}

	r.Code = 200
	r.Data = user
	return
}
