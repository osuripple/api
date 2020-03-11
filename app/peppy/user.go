// Package peppy implements the osu! API as defined on the osu-api repository wiki (https://github.com/ppy/osu-api/wiki).
package peppy

import (
	"database/sql"
	"fmt"
	"strconv"

	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/thehowl/go-osuapi"
	"github.com/valyala/fasthttp"
	"gopkg.in/redis.v5"
	"zxq.co/ripple/ocl"
	"zxq.co/ripple/rippleapi/common"
)

// R is a redis client.
var R *redis.Client

// GetUser retrieves general user information.
func GetUser(c *fasthttp.RequestCtx, db *sqlx.DB) {
	if query(c, "u") == "" {
		json(c, 200, defaultResponse)
		return
	}
	var user osuapi.User
	whereClause, p := genUser(c, db)
	whereClause = "WHERE " + whereClause

	mode := genmode(query(c, "m"))
	isRelax := query(c, "relax") == "1"

	table := "users_stats"
	var classicJoin string
	if isRelax {
		table = "users_stats_relax"
		// 'country' is in users_stats only, we need to join with both
		// users_stats_relax and users_stats 😪
		classicJoin = "LEFT JOIN users_stats USING(id)"
	}

	err := db.QueryRow(fmt.Sprintf(
		`SELECT
			users.id, users.username,
			s.playcount_%s, s.ranked_score_%s, s.total_score_%s,
			s.pp_%s, s.avg_accuracy_%s,
			country
		FROM users
		LEFT JOIN %s AS s USING(id)
		%s
		%s
		LIMIT 1`,
		mode, mode, mode, mode, mode, table, classicJoin, whereClause,
	), p).Scan(
		&user.UserID, &user.Username,
		&user.Playcount, &user.RankedScore, &user.TotalScore,
		&user.PP, &user.Accuracy,
		&user.Country,
	)
	if err != nil {
		json(c, 200, defaultResponse)
		if err != sql.ErrNoRows {
			common.Err(c, err)
		}
		return
	}

	var suffix string
	if isRelax {
		suffix = ":relax"
	}
	user.Rank = int(R.ZRevRank("ripple:leaderboard:"+mode+suffix, strconv.Itoa(user.UserID)).Val()) + 1
	user.CountryRank = int(R.ZRevRank("ripple:leaderboard:"+mode+":"+strings.ToLower(user.Country)+suffix, strconv.Itoa(user.UserID)).Val()) + 1
	user.Level = ocl.GetLevelPrecise(user.TotalScore)

	json(c, 200, []osuapi.User{user})
}
