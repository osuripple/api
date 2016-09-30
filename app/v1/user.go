// Package v1 implements the first version of the Ripple API.
package v1

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"git.zxq.co/ripple/ocl"
	"git.zxq.co/ripple/rippleapi/common"
)

type userData struct {
	ID             int                  `json:"id"`
	Username       string               `json:"username"`
	UsernameAKA    string               `json:"username_aka"`
	RegisteredOn   common.UnixTimestamp `json:"registered_on"`
	Privileges     uint64               `json:"privileges"`
	LatestActivity common.UnixTimestamp `json:"latest_activity"`
	Country        string               `json:"country"`
}

// UsersGET is the API handler for GET /users
func UsersGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return *shouldRet
	}

	query := `
SELECT users.id, users.username, register_datetime, privileges,
	latest_activity, users_stats.username_aka,
	users_stats.country
FROM users
INNER JOIN users_stats
ON users.id=users_stats.id
WHERE ` + whereClause + ` AND ` + md.User.OnlyUserPublic(true) + `
LIMIT 1`
	return userPuts(md, md.DB.QueryRowx(query, param))
}

type userPutsUserData struct {
	common.ResponseBase
	userData
}

func userPuts(md common.MethodData, row *sqlx.Row) common.CodeMessager {
	var err error
	var user userPutsUserData

	err = row.StructScan(&user.userData)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such user was found!")
	case err != nil:
		md.Err(err)
		return Err500
	}

	user.Code = 200
	return user
}

// UserSelfGET is a shortcut for /users/id/self. (/users/self)
func UserSelfGET(md common.MethodData) common.CodeMessager {
	md.C.Request.URL.RawQuery = "id=self&" + md.C.Request.URL.RawQuery
	return UsersGET(md)
}

type whatIDResponse struct {
	common.ResponseBase
	ID int `json:"id"`
}

// UserWhatsTheIDGET is an API request that only returns an user's ID.
func UserWhatsTheIDGET(md common.MethodData) common.CodeMessager {
	var (
		r          whatIDResponse
		privileges uint64
	)
	err := md.DB.QueryRow("SELECT id, privileges FROM users WHERE username = ? LIMIT 1", md.Query("name")).Scan(&r.ID, &privileges)
	if err != nil || ((privileges&uint64(common.UserPrivilegePublic)) == 0 &&
		(md.User.UserPrivileges&common.AdminPrivilegeManageUsers == 0)) {
		return common.SimpleResponse(404, "That user could not be found!")
	}
	r.Code = 200
	return r
}

type modeData struct {
	RankedScore           uint64  `json:"ranked_score"`
	TotalScore            uint64  `json:"total_score"`
	PlayCount             int     `json:"playcount"`
	ReplaysWatched        int     `json:"replays_watched"`
	TotalHits             int     `json:"total_hits"`
	Level                 float64 `json:"level"`
	Accuracy              float64 `json:"accuracy"`
	PP                    int     `json:"pp"`
	GlobalLeaderboardRank *int    `json:"global_leaderboard_rank"`
}
type userFullResponse struct {
	common.ResponseBase
	userData
	STD           modeData      `json:"std"`
	Taiko         modeData      `json:"taiko"`
	CTB           modeData      `json:"ctb"`
	Mania         modeData      `json:"mania"`
	PlayStyle     int           `json:"play_style"`
	FavouriteMode int           `json:"favourite_mode"`
	Badges        []singleBadge `json:"badges"`
	CustomBadge   *singleBadge  `json:"custom_badge"`
}

// UserFullGET gets all of an user's information, with one exception: their userpage.
func UserFullGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users")
	if shouldRet != nil {
		return *shouldRet
	}

	// Hellest query I've ever done.
	query := `
SELECT
	users.id, users.username, users.register_datetime, users.privileges, users.latest_activity,

	users_stats.username_aka, users_stats.country, users_stats.play_style, users_stats.favourite_mode,

	users_stats.custom_badge_icon, users_stats.custom_badge_name, users_stats.can_custom_badge, 
	users_stats.show_custom_badge,

	users_stats.ranked_score_std, users_stats.total_score_std, users_stats.playcount_std,
	users_stats.replays_watched_std, users_stats.total_hits_std,
	users_stats.avg_accuracy_std, users_stats.pp_std, leaderboard_std.position as std_position,

	users_stats.ranked_score_taiko, users_stats.total_score_taiko, users_stats.playcount_taiko,
	users_stats.replays_watched_taiko, users_stats.total_hits_taiko,
	users_stats.avg_accuracy_taiko, users_stats.pp_taiko, leaderboard_taiko.position as taiko_position,

	users_stats.ranked_score_ctb, users_stats.total_score_ctb, users_stats.playcount_ctb,
	users_stats.replays_watched_ctb, users_stats.total_hits_ctb,
	users_stats.avg_accuracy_ctb, users_stats.pp_ctb, leaderboard_ctb.position as ctb_position,

	users_stats.ranked_score_mania, users_stats.total_score_mania, users_stats.playcount_mania,
	users_stats.replays_watched_mania, users_stats.total_hits_mania,
	users_stats.avg_accuracy_mania, users_stats.pp_mania, leaderboard_mania.position as mania_position

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
WHERE ` + whereClause + ` AND ` + md.User.OnlyUserPublic(true) + `
LIMIT 1
`
	// Fuck.
	r := userFullResponse{}
	var (
		b    singleBadge
		can  bool
		show bool
	)
	err := md.DB.QueryRow(query, param).Scan(
		&r.ID, &r.Username, &r.RegisteredOn, &r.Privileges, &r.LatestActivity,

		&r.UsernameAKA, &r.Country,
		&r.PlayStyle, &r.FavouriteMode,

		&b.Icon, &b.Name, &can, &show,

		&r.STD.RankedScore, &r.STD.TotalScore, &r.STD.PlayCount,
		&r.STD.ReplaysWatched, &r.STD.TotalHits,
		&r.STD.Accuracy, &r.STD.PP, &r.STD.GlobalLeaderboardRank,

		&r.Taiko.RankedScore, &r.Taiko.TotalScore, &r.Taiko.PlayCount,
		&r.Taiko.ReplaysWatched, &r.Taiko.TotalHits,
		&r.Taiko.Accuracy, &r.Taiko.PP, &r.Taiko.GlobalLeaderboardRank,

		&r.CTB.RankedScore, &r.CTB.TotalScore, &r.CTB.PlayCount,
		&r.CTB.ReplaysWatched, &r.CTB.TotalHits,
		&r.CTB.Accuracy, &r.CTB.PP, &r.CTB.GlobalLeaderboardRank,

		&r.Mania.RankedScore, &r.Mania.TotalScore, &r.Mania.PlayCount,
		&r.Mania.ReplaysWatched, &r.Mania.TotalHits,
		&r.Mania.Accuracy, &r.Mania.PP, &r.Mania.GlobalLeaderboardRank,
	)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "That user could not be found!")
	case err != nil:
		md.Err(err)
		return Err500
	}

	can = can && show && common.UserPrivileges(r.Privileges)&common.UserPrivilegeDonor > 0
	if can && (b.Name != "" || b.Icon != "") {
		r.CustomBadge = &b
	}

	for _, m := range []*modeData{&r.STD, &r.Taiko, &r.CTB, &r.Mania} {
		m.Level = ocl.GetLevelPrecise(int64(m.TotalScore))
	}

	rows, err := md.DB.Query("SELECT b.id, b.name, b.icon FROM user_badges ub "+
		"LEFT JOIN badges b ON ub.badge = b.id WHERE user = ?", r.ID)
	if err != nil {
		md.Err(err)
	}

	for rows.Next() {
		var badge singleBadge
		err := rows.Scan(&badge.ID, &badge.Name, &badge.Icon)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Badges = append(r.Badges, badge)
	}

	r.Code = 200
	return r
}

type userpageResponse struct {
	common.ResponseBase
	Userpage *string `json:"userpage"`
}

// UserUserpageGET gets an user's userpage, as in the customisable thing.
func UserUserpageGET(md common.MethodData) common.CodeMessager {
	shouldRet, whereClause, param := whereClauseUser(md, "users_stats")
	if shouldRet != nil {
		return *shouldRet
	}
	var r userpageResponse
	err := md.DB.QueryRow("SELECT userpage_content FROM users_stats WHERE "+whereClause+" LIMIT 1", param).Scan(&r.Userpage)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such user!")
	case err != nil:
		md.Err(err)
		return Err500
	}
	r.Code = 200
	return r
}

func whereClauseUser(md common.MethodData, tableName string) (*common.CodeMessager, string, interface{}) {
	switch {
	case md.Query("id") == "self":
		return nil, tableName + ".id = ?", md.ID()
	case md.Query("id") != "":
		id, err := strconv.Atoi(md.Query("id"))
		if err != nil {
			a := common.SimpleResponse(400, "please pass a valid user ID")
			return &a, "", nil
		}
		return nil, tableName + ".id = ?", id
	case md.Query("name") != "":
		return nil, tableName + ".username = ?", md.Query("name")
	}
	a := common.SimpleResponse(400, "you need to pass either querystring parameters name or id")
	return &a, "", nil
}

type userLookupResponse struct {
	common.ResponseBase
	Users []lookupUser `json:"users"`
}
type lookupUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// UserLookupGET does a quick lookup of users beginning with the passed
// querystring value name.
func UserLookupGET(md common.MethodData) common.CodeMessager {
	name := strings.NewReplacer(
		"%", "\\%",
		"_", "\\_",
		"\\", "\\\\",
	).Replace(md.Query("name"))
	if name == "" {
		return common.SimpleResponse(400, "please provide an username to start searching")
	}
	name = "%" + name + "%"
	rows, err := md.DB.Query("SELECT users.id, users.username FROM users WHERE username LIKE ? AND "+
		md.User.OnlyUserPublic(true)+" LIMIT 25", name)
	if err != nil {
		md.Err(err)
		return Err500
	}

	var r userLookupResponse

	for rows.Next() {
		var l lookupUser
		err := rows.Scan(&l.ID, &l.Username)
		if err != nil {
			continue // can't be bothered to handle properly
		}
		r.Users = append(r.Users, l)
	}

	r.Code = 200
	return r
}
