package app

import (
	"fmt"

	"git.zxq.co/ripple/rippleapi/app/internals"
	"git.zxq.co/ripple/rippleapi/app/peppy"
	"git.zxq.co/ripple/rippleapi/app/v1"
	"git.zxq.co/ripple/rippleapi/common"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
	cf common.Conf
)

// Start begins taking HTTP connections.
func Start(conf common.Conf, dbO *sqlx.DB) *gin.Engine {
	db = dbO
	cf = conf

	setUpLimiter()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	if conf.SentryDSN != "" {
		ravenClient, err := raven.New(conf.SentryDSN)
		ravenClient.SetRelease(common.Version)
		if err != nil {
			fmt.Println(err)
		} else {
			r.Use(Recovery(ravenClient, false))
		}
	}

	api := r.Group("/api")
	{
		gv1 := api.Group("/v1")
		{
			gv1.POST("/tokens", Method(v1.TokenNewPOST))
			gv1.POST("/tokens/new", Method(v1.TokenNewPOST))
			gv1.GET("/tokens/self/delete", Method(v1.TokenSelfDeleteGET))

			// Auth-free API endpoints (public data)
			gv1.GET("/ping", Method(v1.PingGET))
			gv1.GET("/surprise_me", Method(v1.SurpriseMeGET))
			gv1.GET("/privileges", Method(v1.PrivilegesGET))
			gv1.GET("/doc", Method(v1.DocGET))
			gv1.GET("/doc/content", Method(v1.DocContentGET))
			gv1.GET("/doc/rules", Method(v1.DocRulesGET))
			gv1.GET("/users", Method(v1.UsersGET))
			gv1.GET("/users/whatid", Method(v1.UserWhatsTheIDGET))
			gv1.GET("/users/full", Method(v1.UserFullGET))
			gv1.GET("/users/userpage", Method(v1.UserUserpageGET))
			gv1.GET("/users/lookup", Method(v1.UserLookupGET))
			gv1.GET("/users/scores/best", Method(v1.UserScoresBestGET))
			gv1.GET("/users/scores/recent", Method(v1.UserScoresRecentGET))
			gv1.GET("/badges", Method(v1.BadgesGET))
			gv1.GET("/beatmaps", Method(v1.BeatmapGET))
			gv1.GET("/leaderboard", Method(v1.LeaderboardGET))
			gv1.GET("/tokens", Method(v1.TokenGET))
			gv1.GET("/users/self", Method(v1.UserSelfGET))
			gv1.GET("/tokens/self", Method(v1.TokenSelfGET))
			gv1.GET("/blog/posts", Method(v1.BlogPostsGET))
			gv1.GET("/blog/posts/content", Method(v1.BlogPostsContentGET))

			// ReadConfidential privilege required
			gv1.GET("/friends", Method(v1.FriendsGET, common.PrivilegeReadConfidential))
			gv1.GET("/friends/with", Method(v1.FriendsWithGET, common.PrivilegeReadConfidential))

			// Write privilege required
			gv1.GET("/friends/add", Method(v1.FriendsAddGET, common.PrivilegeWrite))
			gv1.GET("/friends/del", Method(v1.FriendsDelGET, common.PrivilegeWrite))

			// Admin: beatmap
			gv1.POST("/beatmaps/set_status", Method(v1.BeatmapSetStatusPOST, common.PrivilegeBeatmap))
			gv1.GET("/beatmaps/ranked_frozen_full", Method(v1.BeatmapRankedFrozenFullGET, common.PrivilegeBeatmap))

			// Admin: user managing
			gv1.POST("/users/manage/set_allowed", Method(v1.UserManageSetAllowedPOST, common.PrivilegeManageUser))

			// M E T A
			// E     T    "wow thats so meta"
			// T     E                  -- the one who said "wow thats so meta"
			// A T E M
			gv1.GET("/meta/restart", Method(v1.MetaRestartGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/kill", Method(v1.MetaKillGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/up_since", Method(v1.MetaUpSinceGET, common.PrivilegeAPIMeta))
			gv1.GET("/meta/update", Method(v1.MetaUpdateGET, common.PrivilegeAPIMeta))

			// User Managing + meta
			gv1.GET("/tokens/fix_privileges", Method(v1.TokenFixPrivilegesGET,
				common.PrivilegeManageUser, common.PrivilegeAPIMeta))
		}

		api.GET("/status", internals.Status)

		// peppyapi
		api.GET("/get_user", PeppyMethod(peppy.GetUser))
		api.GET("/get_match", PeppyMethod(peppy.GetMatch))
		api.GET("/get_user_recent", PeppyMethod(peppy.GetUserRecent))
		api.GET("/get_user_best", PeppyMethod(peppy.GetUserBest))
	}

	r.NoRoute(v1.Handle404)

	return r
	/*if conf.Unix {
		panic(r.RunUnix(conf.ListenTo))
	}
	panic(r.Run(conf.ListenTo))*/
}
