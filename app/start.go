package app

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	fhr "github.com/buaazp/fasthttprouter"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/sqlx"
	"github.com/serenize/snaker"
	"gopkg.in/redis.v5"
	"zxq.co/ripple/rippleapi/app/internals"
	"zxq.co/ripple/rippleapi/app/peppy"
	"zxq.co/ripple/rippleapi/app/v1"
	"zxq.co/ripple/rippleapi/common"
)

var (
	db    *sqlx.DB
	cf    common.Conf
	doggo *statsd.Client
	red   *redis.Client
)

var commonClusterfucks = map[string]string{
	"RegisteredOn": "register_datetime",
	"UsernameAKA":  "username_aka",
}

// Start begins taking HTTP connections.
func Start(conf common.Conf, dbO *sqlx.DB) *fhr.Router {
	db = dbO
	cf = conf

	db.MapperFunc(func(s string) string {
		if x, ok := commonClusterfucks[s]; ok {
			return x
		}
		return snaker.CamelToSnake(s)
	})

	r := fhr.New()
	// TODO: add back gzip
	// TODO: add logging
	// TODO: add sentry panic recovering

	// sentry
	if conf.SentryDSN != "" {
		ravenClient, err := raven.New(conf.SentryDSN)
		ravenClient.SetRelease(common.Version)
		if err != nil {
			fmt.Println(err)
		} else {
			// r.Use(Recovery(ravenClient, false))
			common.RavenClient = ravenClient
		}
	}

	// datadog
	var err error
	doggo, err = statsd.New("127.0.0.1:8125")
	if err != nil {
		fmt.Println(err)
	}
	doggo.Namespace = "api."
	// r.Use(func(c *gin.Context) {
	// 	doggo.Incr("requests", nil, 1)
	// })

	// redis
	red = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	// token updater
	go tokenUpdater(db)

	// peppyapi
	{
		r.GET("/api/get_user", PeppyMethod(peppy.GetUser))
		r.GET("/api/get_match", PeppyMethod(peppy.GetMatch))
		r.GET("/api/get_user_recent", PeppyMethod(peppy.GetUserRecent))
		r.GET("/api/get_user_best", PeppyMethod(peppy.GetUserBest))
		r.GET("/api/get_scores", PeppyMethod(peppy.GetScores))
		r.GET("/api/get_beatmaps", PeppyMethod(peppy.GetBeatmap))
	}

	// v1 API
	{
		r.POST("/api/v1/tokens", Method(v1.TokenNewPOST))
		r.POST("/api/v1/tokens/new", Method(v1.TokenNewPOST))
		r.POST("/api/v1/tokens/self/delete", Method(v1.TokenSelfDeletePOST))

		// Auth-free API endpoints (public data)
		r.GET("/api/v1/ping", Method(v1.PingGET))
		r.GET("/api/v1/surprise_me", Method(v1.SurpriseMeGET))
		r.GET("/api/v1/doc", Method(v1.DocGET))
		r.GET("/api/v1/doc/content", Method(v1.DocContentGET))
		r.GET("/api/v1/doc/rules", Method(v1.DocRulesGET))
		r.GET("/api/v1/users", Method(v1.UsersGET))
		r.GET("/api/v1/users/whatid", Method(v1.UserWhatsTheIDGET))
		r.GET("/api/v1/users/full", Method(v1.UserFullGET))
		r.GET("/api/v1/users/userpage", Method(v1.UserUserpageGET))
		r.GET("/api/v1/users/lookup", Method(v1.UserLookupGET))
		r.GET("/api/v1/users/scores/best", Method(v1.UserScoresBestGET))
		r.GET("/api/v1/users/scores/recent", Method(v1.UserScoresRecentGET))
		r.GET("/api/v1/badges", Method(v1.BadgesGET))
		r.GET("/api/v1/beatmaps", Method(v1.BeatmapGET))
		r.GET("/api/v1/leaderboard", Method(v1.LeaderboardGET))
		r.GET("/api/v1/tokens", Method(v1.TokenGET))
		r.GET("/api/v1/users/self", Method(v1.UserSelfGET))
		r.GET("/api/v1/tokens/self", Method(v1.TokenSelfGET))
		r.GET("/api/v1/blog/posts", Method(v1.BlogPostsGET))
		r.GET("/api/v1/scores", Method(v1.ScoresGET))
		r.GET("/api/v1/beatmaps/rank_requests/status", Method(v1.BeatmapRankRequestsStatusGET))

		// ReadConfidential privilege required
		r.GET("/api/v1/friends", Method(v1.FriendsGET, common.PrivilegeReadConfidential))
		r.GET("/api/v1/friends/with", Method(v1.FriendsWithGET, common.PrivilegeReadConfidential))
		r.GET("/api/v1/users/self/donor_info", Method(v1.UsersSelfDonorInfoGET, common.PrivilegeReadConfidential))
		r.GET("/api/v1/users/self/favourite_mode", Method(v1.UsersSelfFavouriteModeGET, common.PrivilegeReadConfidential))
		r.GET("/api/v1/users/self/settings", Method(v1.UsersSelfSettingsGET, common.PrivilegeReadConfidential))

		// Write privilege required
		r.POST("/api/v1/friends/add", Method(v1.FriendsAddPOST, common.PrivilegeWrite))
		r.POST("/api/v1/friends/del", Method(v1.FriendsDelPOST, common.PrivilegeWrite))
		r.POST("/api/v1/users/self/settings", Method(v1.UsersSelfSettingsPOST, common.PrivilegeWrite))
		r.POST("/api/v1/users/self/userpage", Method(v1.UserSelfUserpagePOST, common.PrivilegeWrite))
		r.POST("/api/v1/beatmaps/rank_requests", Method(v1.BeatmapRankRequestsSubmitPOST, common.PrivilegeWrite))

		// Admin: beatmap
		r.POST("/api/v1/beatmaps/set_status", Method(v1.BeatmapSetStatusPOST, common.PrivilegeBeatmap))
		r.GET("/api/v1/beatmaps/ranked_frozen_full", Method(v1.BeatmapRankedFrozenFullGET, common.PrivilegeBeatmap))

		// Admin: user managing
		r.POST("/api/v1/users/manage/set_allowed", Method(v1.UserManageSetAllowedPOST, common.PrivilegeManageUser))

		// M E T A
		// E     T    "wow thats so meta"
		// T     E                  -- the one who said "wow thats so meta"
		// A T E M
		r.GET("/api/v1/meta/restart", Method(v1.MetaRestartGET, common.PrivilegeAPIMeta))
		r.GET("/api/v1/meta/kill", Method(v1.MetaKillGET, common.PrivilegeAPIMeta))
		r.GET("/api/v1/meta/up_since", Method(v1.MetaUpSinceGET, common.PrivilegeAPIMeta))
		r.GET("/api/v1/meta/update", Method(v1.MetaUpdateGET, common.PrivilegeAPIMeta))

		// User Managing + meta
		r.POST("/api/v1/tokens/fix_privileges", Method(v1.TokenFixPrivilegesPOST,
			common.PrivilegeManageUser, common.PrivilegeAPIMeta))
	}

	// in the new osu-web, the old endpoints are also in /v1 it seems. So /shrug
	{
		r.GET("/api/v1/get_user", PeppyMethod(peppy.GetUser))
		r.GET("/api/v1/get_match", PeppyMethod(peppy.GetMatch))
		r.GET("/api/v1/get_user_recent", PeppyMethod(peppy.GetUserRecent))
		r.GET("/api/v1/get_user_best", PeppyMethod(peppy.GetUserBest))
		r.GET("/api/v1/get_scores", PeppyMethod(peppy.GetScores))
		r.GET("/api/v1/get_beatmaps", PeppyMethod(peppy.GetBeatmap))
	}

	r.GET("/api/status", internals.Status)

	r.NotFound = v1.Handle404

	return r
}
