package app

import (
	"database/sql"

	"git.zxq.co/ripple/rippleapi/app/internals"
	"git.zxq.co/ripple/rippleapi/app/peppy"
	"git.zxq.co/ripple/rippleapi/app/v1"
	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

// Start begins taking HTTP connections.
func Start(conf common.Conf, dbO *sql.DB) *gin.Engine {
	db = dbO
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression), ErrorHandler())

	api := r.Group("/api")
	{
		gv1 := api.Group("/v1")
		{
			gv1.POST("/tokens/new", Method(v1.TokenNewPOST))
			gv1.GET("/tokens/self/delete", Method(v1.TokenSelfDeleteGET))

			// Auth-free API endpoints
			gv1.GET("/ping", Method(v1.PingGET))
			gv1.GET("/surprise_me", Method(v1.SurpriseMeGET))
			gv1.GET("/privileges", Method(v1.PrivilegesGET))
			gv1.GET("/doc", Method(v1.DocGET))
			gv1.GET("/doc/content", Method(v1.DocContentGET))
			gv1.GET("/doc/rules", Method(v1.DocRulesGET))

			// Read privilege required
			gv1.GET("/users", Method(v1.UsersGET, common.PrivilegeRead))
			gv1.GET("/users/self", Method(v1.UserSelfGET, common.PrivilegeRead))
			gv1.GET("/users/whatid", Method(v1.UserWhatsTheIDGET, common.PrivilegeRead))
			gv1.GET("/users/full", Method(v1.UserFullGET, common.PrivilegeRead))
			gv1.GET("/users/userpage", Method(v1.UserUserpageGET, common.PrivilegeRead))
			gv1.GET("/users/lookup", Method(v1.UserLookupGET, common.PrivilegeRead))
			gv1.GET("/users/scores/best", Method(v1.UserScoresBestGET, common.PrivilegeRead))
			gv1.GET("/users/scores/recent", Method(v1.UserScoresRecentGET, common.PrivilegeRead))
			gv1.GET("/badges", Method(v1.BadgesGET, common.PrivilegeRead))
			gv1.GET("/beatmaps", Method(v1.BeatmapGET, common.PrivilegeRead))
			gv1.GET("/leaderboard", Method(v1.LeaderboardGET, common.PrivilegeRead))
			gv1.GET("/tokens", Method(v1.TokenGET, common.PrivilegeRead))
			gv1.GET("/tokens/self", Method(v1.TokenSelfGET, common.PrivilegeRead))

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
