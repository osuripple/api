package app

import (
	"database/sql"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/osuripple/api/app/v1"
	"github.com/osuripple/api/common"
)

// Start begins taking HTTP connections.
func Start(conf common.Conf, db *sql.DB) {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	api := r.Group("/api")
	{
		gv1 := api.Group("/v1")
		{
			// Auth-free API endpoints
			gv1.GET("/ping", Method(v1.Ping, db))
			gv1.GET("/surprise_me", Method(v1.SurpriseMe, db))
			gv1.GET("/privileges", Method(v1.PrivilegesGET, db))

			// Read privilege required
			gv1.GET("/users/:id", Method(v1.UserGET, db, common.PrivilegeRead))
			gv1.GET("/badges", Method(v1.BadgesGET, db, common.PrivilegeRead))
			gv1.GET("/badges/:id", Method(v1.BadgeByIDGET, db, common.PrivilegeRead))
		}
	}

	r.NoRoute(v1.Handle404)
	if conf.Unix {
		panic(r.RunUnix(conf.ListenTo))
	}
	panic(r.Run(conf.ListenTo))
}
