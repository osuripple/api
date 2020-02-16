package v1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	redis "gopkg.in/redis.v5"

	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

type leaderboardUser struct {
	userData
	ChosenMode    modeData `json:"chosen_mode"`
	PlayStyle     int      `json:"play_style"`
	FavouriteMode int      `json:"favourite_mode"`
}

type leaderboardResponse struct {
	common.ResponseBase
	Users []leaderboardUser `json:"users"`
}

const lbUserQuery = `
SELECT
	users.id, users.username, users.register_datetime, users.privileges, users.latest_activity,

	full_stats.username_aka, full_stats.country,
	full_stats.play_style, full_stats.favourite_mode,

	stats.ranked_score_%[1]s, stats.total_score_%[1]s, stats.playcount_%[1]s,
	full_stats.replays_watched_%[1]s, stats.total_hits_%[1]s,
	stats.avg_accuracy_%[1]s, stats.pp_%[1]s
FROM users
INNER JOIN %[2]s AS stats USING(id) JOIN users_stats AS full_stats USING(id)
WHERE users.id IN (?)
`

// LeaderboardGET gets the leaderboard.
func LeaderboardGET(md common.MethodData) common.CodeMessager {
	m := getMode(md.Query("mode"))

	// md.Query.Country
	p := common.Int(md.Query("p")) - 1
	if p < 0 {
		p = 0
	}
	l := common.InString(1, md.Query("l"), 500, 50)

	key := "ripple:leaderboard:" + m
	if md.Query("country") != "" {
		key += ":" + md.Query("country")
	}
	isRelax := md.Query("relax") == "1"
	var table string
	if isRelax {
		key += ":relax"
		table = "users_stats_relax"
	} else {
		table = "users_stats"
	}

	results, err := md.R.ZRevRange(key, int64(p*l), int64(p*l+l-1)).Result()
	if err != nil {
		md.Err(err)
		return Err500
	}

	var resp leaderboardResponse
	resp.Code = 200

	if len(results) == 0 {
		return resp
	}

	query := fmt.Sprintf(lbUserQuery+` ORDER BY stats.pp_%[1]s DESC, stats.ranked_score_%[1]s DESC`, m, table)
	query, params, _ := sqlx.In(query, results)
	rows, err := md.DB.Query(query, params...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	for rows.Next() {
		var u leaderboardUser
		err := rows.Scan(
			&u.ID, &u.Username, &u.RegisteredOn, &u.Privileges, &u.LatestActivity,

			&u.UsernameAKA, &u.Country, &u.PlayStyle, &u.FavouriteMode,

			&u.ChosenMode.RankedScore, &u.ChosenMode.TotalScore, &u.ChosenMode.PlayCount,
			&u.ChosenMode.ReplaysWatched, &u.ChosenMode.TotalHits,
			&u.ChosenMode.Accuracy, &u.ChosenMode.PP,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		u.ChosenMode.Level = ocl.GetLevelPrecise(int64(u.ChosenMode.TotalScore))
		if i := leaderboardPosition(md.R, m, u.ID, isRelax); i != nil {
			u.ChosenMode.GlobalLeaderboardRank = i
		}
		if i := countryPosition(md.R, m, u.ID, u.Country, isRelax); i != nil {
			u.ChosenMode.CountryLeaderboardRank = i
		}
		resp.Users = append(resp.Users, u)
	}
	return resp
}

func leaderboardPosition(r *redis.Client, mode string, user int, relax bool) *int {
	var suffix string
	if relax {
		suffix = ":relax"
	} else {
		suffix = ""
	}
	return _position(r, "ripple:leaderboard:"+mode+suffix, user)
}

func countryPosition(r *redis.Client, mode string, user int, country string, relax bool) *int {
	var suffix string
	if relax {
		suffix = ":relax"
	} else {
		suffix = ""
	}
	return _position(r, "ripple:leaderboard:"+mode+":"+strings.ToLower(country)+suffix, user)
}

func _position(r *redis.Client, key string, user int) *int {
	res := r.ZRevRank(key, strconv.Itoa(user))
	if res.Err() == redis.Nil {
		return nil
	}
	x := int(res.Val()) + 1
	return &x
}
