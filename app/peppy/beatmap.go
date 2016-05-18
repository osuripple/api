package peppy

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// GetBeatmap retrieves general beatmap information.
func GetBeatmap(c *gin.Context, db *sql.DB) {
	var whereClauses []string
	var params []string

	// since value is not stored, silently ignore
	if c.Query("s") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmapset_id = ?")
		params = append(params, c.Query("s"))
	}
	if c.Query("b") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmap_id = ?")
		params = append(params, c.Query("b"))
	}
	if c.Query("u") != "" {
		wc, p := genUser(c, db)
		whereClauses = append(whereClauses, wc)
		params = append(params, p)
	}
	// silently ignore m
	// silently ignore a
	if c.Query("h") != "" {
		whereClauses = append(whereClauses, "beatmaps.beatmap_md5 = ?")
		params = append(params, c.Query("h"))
	}

	//bm := osuapi.Beatmap{}

	//db.Query("SELECT beatmaps.beatmapset_id, beatmaps.beatmap FROM ")
}
