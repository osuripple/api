package v1

import (
	"fmt"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
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
	users.id, users.username, users.register_datetime, users.rank, users.latest_activity,
	
	users_stats.username_aka, users_stats.country, users_stats.show_country,
	users_stats.play_style, users_stats.favourite_mode,
	
	users_stats.ranked_score_%[1]s, users_stats.total_score_%[1]s, users_stats.playcount_%[1]s,
	users_stats.replays_watched_%[1]s, users_stats.total_hits_%[1]s,
	users_stats.avg_accuracy_%[1]s, users_stats.pp_%[1]s, leaderboard_%[1]s.position as %[1]s_position
FROM leaderboard_%[1]s
INNER JOIN users ON users.id = leaderboard_%[1]s.user
INNER JOIN users_stats ON users_stats.id = leaderboard_%[1]s.user
%[2]s`

// LeaderboardGET gets the leaderboard.
func LeaderboardGET(md common.MethodData) common.CodeMessager {
	m := getMode(md.C.Query("mode"))
	query := fmt.Sprintf(lbUserQuery, m, `WHERE users.allowed = '1' ORDER BY leaderboard_`+m+`.position `+
		common.Paginate(md.C.Query("p"), md.C.Query("l"), 100))
	rows, err := md.DB.Query(query)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var resp leaderboardResponse
	for rows.Next() {
		var (
			u              leaderboardUser
			register       int64
			latestActivity int64
			showCountry    bool
		)
		err := rows.Scan(
			&u.ID, &u.Username, &register, &u.Rank, &latestActivity,

			&u.UsernameAKA, &u.Country, &showCountry,
			&u.PlayStyle, &u.FavouriteMode,

			&u.ChosenMode.RankedScore, &u.ChosenMode.TotalScore, &u.ChosenMode.PlayCount,
			&u.ChosenMode.ReplaysWatched, &u.ChosenMode.TotalHits,
			&u.ChosenMode.Accuracy, &u.ChosenMode.PP, &u.ChosenMode.GlobalLeaderboardRank,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		if !showCountry {
			u.Country = "XX"
		}
		u.RegisteredOn = time.Unix(register, 0)
		u.LatestActivity = time.Unix(latestActivity, 0)
		resp.Users = append(resp.Users, u)
	}
	resp.Code = 200
	return resp
}
