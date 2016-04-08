// Package v1 implements the first version of the Ripple API.
package v1

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		uid = md.ID()
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

	var (
		registeredOn   int64
		latestActivity int64
		badges         string
		showCountry    bool
	)
	err = row.Scan(&user.ID, &user.Username, &registeredOn, &user.Rank, &latestActivity, &user.UsernameAKA, &badges, &user.Country, &showCountry)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "No such user was found!"
		return
	case err != nil:
		md.Err(err)
		r = Err500
		return
	}

	user.RegisteredOn = time.Unix(registeredOn, 0)
	user.LatestActivity = time.Unix(latestActivity, 0)

	user.Badges = badgesToArray(badges)

	user.Country = genCountry(md, user.ID, showCountry, user.Country)

	r.Code = 200
	r.Data = user
	return
}

func badgesToArray(badges string) []int {
	var end []int
	badgesSl := strings.Split(badges, ",")
	for _, badge := range badgesSl {
		if badge != "" && badge != "0" {
			// We are ignoring errors because who really gives a shit if something's gone wrong on our end in this
			// particular thing, we can just silently ignore this.
			nb, err := strconv.Atoi(badge)
			if err == nil && nb != 0 {
				end = append(end, nb)
			}
		}
	}
	return end
}

func genCountry(md common.MethodData, uid int, showCountry bool, country string) string {
	// If the user wants to stay anonymous, don't show their country.
	// This can be overriden if we have the ReadConfidential privilege and the user we are accessing is the token owner.
	if showCountry || (md.User.Privileges.HasPrivilegeReadConfidential() && uid == md.ID()) {
		return country
	}
	return "XX"
}

// UserSelfGET is a shortcut for /users/id/self. (/users/self)
func UserSelfGET(md common.MethodData) common.Response {
	md.C.Params = append(md.C.Params, gin.Param{
		Key:   "id",
		Value: "self",
	})
	return UserByIDGET(md)
}

// UserWhatsTheIDGET is an API request that only returns an user's ID.
func UserWhatsTheIDGET(md common.MethodData) common.Response {
	var (
		id      int
		allowed int
	)
	err := md.DB.QueryRow("SELECT id, allowed FROM users WHERE username = ? LIMIT 1", md.C.Param("username")).Scan(&id, &allowed)
	if err != nil || allowed != 1 {
		return common.Response{
			Code:    404,
			Message: "That user could not be found!",
		}
	}
	return common.Response{
		Code: 200,
		Data: id,
	}
}

type modeData struct {
	RankedScore           uint64  `json:"ranked_score"`
	TotalScore            uint64  `json:"total_score"`
	PlayCount             int     `json:"playcount"`
	ReplaysWatched        int     `json:"replays_watched"`
	TotalHits             int     `json:"total_hits"`
	Level                 float64 `json:"level"`
	Accuracy              float64 `json:"accuracy"`
	GlobalLeaderboardRank int     `json:"global_leaderboard_rank"`
}
type userFullData struct {
	userData
	STD           modeData `json:"std"`
	Taiko         modeData `json:"taiko"`
	CTB           modeData `json:"ctb"`
	Mania         modeData `json:"mania"`
	PlayStyle     int      `json:"play_style"`
	FavouriteMode int      `json:"favourite_mode"`
}

// UserFullGET gets all of an user's information, with one exception: their userpage.
func UserFullGET(md common.MethodData) (r common.Response) {
	// Hellest query I've ever done.
	query := `
SELECT
	users.id, users.username, users.register_datetime, users.rank, users.latest_activity,
	
	users_stats.username_aka, users_stats.badges_shown, users_stats.country, users_stats.show_country,
	users_stats.play_style, users_stats.favourite_mode,
	
	users_stats.ranked_score_std, users_stats.total_score_std, users_stats.playcount_std,
	users_stats.replays_watched_std, users_stats.total_hits_std, users_stats.level_std,
	users_stats.avg_accuracy_std, leaderboard_std.position as std_position,
	
	users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.playcount_taiko,
	users_stats.replays_watched_taiko, users_stats.total_hits_taiko, users_stats.level_taiko,
	users_stats.avg_accuracy_taiko, leaderboard_taiko.position as taiko_position,

	users_stats.ranked_score_ctb, users_stats.total_score_ctb, users_stats.playcount_ctb,
	users_stats.replays_watched_ctb, users_stats.total_hits_ctb, users_stats.level_ctb,
	users_stats.avg_accuracy_ctb, leaderboard_ctb.position as ctb_position,

	users_stats.ranked_score_mania, users_stats.total_score_mania, users_stats.playcount_mania,
	users_stats.replays_watched_mania, users_stats.total_hits_mania, users_stats.level_mania,
	users_stats.avg_accuracy_mania, leaderboard_mania.position as mania_position

FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
LEFT JOIN leaderboard_std
ON users.id=leaderboard_std.user
LEFT JOIN leaderboard_taiko
ON users.id=leaderboard_taiko.user
LEFT JOIN leaderboard_ctb
ON users.id=leaderboard_ctb.user
LEFT JOIN leaderboard_mania
ON users.id=leaderboard_mania.user
WHERE users.id=?
LIMIT 1
`
	// Fuck.
	fd := userFullData{}
	var (
		badges         string
		country        string
		showCountry    bool
		registeredOn   int64
		latestActivity int64
	)
	err := md.DB.QueryRow(query, md.C.Param("id")).Scan(
		&fd.ID, &fd.Username, &registeredOn, &fd.Rank, &latestActivity,

		&fd.UsernameAKA, &badges, &country, &showCountry,
		&fd.PlayStyle, &fd.FavouriteMode,

		&fd.STD.RankedScore, &fd.STD.TotalScore, &fd.STD.PlayCount,
		&fd.STD.ReplaysWatched, &fd.STD.TotalHits, &fd.STD.Level,
		&fd.STD.Accuracy, &fd.STD.GlobalLeaderboardRank,

		&fd.Taiko.RankedScore, &fd.Taiko.TotalScore, &fd.Taiko.PlayCount,
		&fd.Taiko.ReplaysWatched, &fd.Taiko.TotalHits, &fd.Taiko.Level,
		&fd.Taiko.Accuracy, &fd.Taiko.GlobalLeaderboardRank,

		&fd.CTB.RankedScore, &fd.CTB.TotalScore, &fd.CTB.PlayCount,
		&fd.CTB.ReplaysWatched, &fd.CTB.TotalHits, &fd.CTB.Level,
		&fd.CTB.Accuracy, &fd.CTB.GlobalLeaderboardRank,

		&fd.Mania.RankedScore, &fd.Mania.TotalScore, &fd.Mania.PlayCount,
		&fd.Mania.ReplaysWatched, &fd.Mania.TotalHits, &fd.Mania.Level,
		&fd.Mania.Accuracy, &fd.Mania.GlobalLeaderboardRank,
	)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "That user could not be found!"
		return
	case err != nil:
		md.Err(err)
		r = Err500
		return
	}

	fd.Country = genCountry(md, fd.ID, showCountry, country)
	fd.Badges = badgesToArray(badges)

	fd.RegisteredOn = time.Unix(registeredOn, 0)
	fd.LatestActivity = time.Unix(latestActivity, 0)

	r.Code = 200
	r.Data = fd
	return
}

// UserUserpageGET gets an user's userpage, as in the customisable thing.
func UserUserpageGET(md common.MethodData) (r common.Response) {
	var userpage string
	err := md.DB.QueryRow("SELECT userpage_content FROM users_stats WHERE id = ? LIMIT 1", md.C.Param("id")).Scan(&userpage)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "No user with that user ID!"
	case err != nil:
		md.Err(err)
		r = Err500
		return
	}
	r.Code = 200
	r.Data = userpage
	return
}
