package app

import (
	"database/sql"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/osuripple/api/app/v1"
	"github.com/osuripple/api/common"
)

// Start begins taking HTTP connections.
func Start(conf common.Conf, db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression), ErrorHandler())

	api := r.Group("/api")
	{
		gv1 := api.Group("/v1")
		{
			gv1.POST("/token/new", Method(v1.TokenNewPOST, db))

			// Auth-free API endpoints
			gv1.GET("/ping", Method(v1.PingGET, db))
			gv1.GET("/surprise_me", Method(v1.SurpriseMeGET, db))
			gv1.GET("/privileges", Method(v1.PrivilegesGET, db))

			// Read privilege required
			gv1.GET("/users/id/:id", Method(v1.UserByIDGET, db, common.PrivilegeRead))
			gv1.GET("/users/name/:name", Method(v1.UserByNameGET, db, common.PrivilegeRead))
			gv1.GET("/users/self", Method(v1.UserSelfGET, db, common.PrivilegeRead))
			gv1.GET("/users/whatid/:username", Method(v1.UserWhatsTheIDGET, db, common.PrivilegeRead))
			gv1.GET("/users/full/:id", Method(v1.UserFullGET, db, common.PrivilegeRead))
			gv1.GET("/users/userpage/:id", Method(v1.UserUserpageGET, db, common.PrivilegeRead))
			gv1.GET("/badges", Method(v1.BadgesGET, db, common.PrivilegeRead))
			gv1.GET("/badges/:id", Method(v1.BadgeByIDGET, db, common.PrivilegeRead))

			// ReadConfidential privilege required
			gv1.GET("/friends", Method(v1.FriendsGET, db, common.PrivilegeReadConfidential))
			gv1.GET("/friends/with/:id", Method(v1.FriendsWithGET, db, common.PrivilegeReadConfidential))

			// Write privilege required
			gv1.POST("/friends/add", Method(v1.FriendsAddPOST, db, common.PrivilegeWrite))
			gv1.GET("/friends/add/:id", Method(v1.FriendsAddGET, db, common.PrivilegeWrite))
			gv1.POST("/friends/del", Method(v1.FriendsDelPOST, db, common.PrivilegeWrite))
			gv1.GET("/friends/del/:id", Method(v1.FriendsDelGET, db, common.PrivilegeWrite))

			// M E T A
			// E     T    "wow thats so meta"
			// T     E                  -- the one who said "wow that's so meta"
			// A T E M
			gv1.GET("/meta/restart", Method(v1.MetaRestartGET, db, common.PrivilegeAPIMeta))
		}
	}

	r.NoRoute(v1.Handle404)

	return r
	/*if conf.Unix {
		panic(r.RunUnix(conf.ListenTo))
	}
	panic(r.Run(conf.ListenTo))*/
}
