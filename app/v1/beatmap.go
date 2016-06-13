package v1

import (
	"database/sql"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
)

type beatmap struct {
	BeatmapID          int       `json:"beatmap_id"`
	BeatmapsetID       int       `json:"beatmapset_id"`
	BeatmapMD5         string    `json:"beatmap_md5"`
	SongName           string    `json:"song_name"`
	AR                 float32   `json:"ar"`
	OD                 float32   `json:"od"`
	Difficulty         float64   `json:"difficulty"`
	MaxCombo           int       `json:"max_combo"`
	HitLength          int       `json:"hit_length"`
	Ranked             int       `json:"ranked"`
	RankedStatusFrozen int       `json:"ranked_status_frozen"`
	LatestUpdate       time.Time `json:"latest_update"`
}

type beatmapMayOrMayNotExist struct {
	BeatmapID          *int
	BeatmapsetID       *int
	BeatmapMD5         *string
	SongName           *string
	AR                 *float32
	OD                 *float32
	Difficulty         *float64
	MaxCombo           *int
	HitLength          *int
	Ranked             *int
	RankedStatusFrozen *int
	LatestUpdate       *time.Time
}

func (b *beatmapMayOrMayNotExist) toBeatmap() *beatmap {
	if b == nil || b.BeatmapID == nil {
		return nil
	}
	return &beatmap{
		BeatmapID:          *b.BeatmapID,
		BeatmapsetID:       *b.BeatmapsetID,
		BeatmapMD5:         *b.BeatmapMD5,
		SongName:           *b.SongName,
		AR:                 *b.AR,
		OD:                 *b.OD,
		Difficulty:         *b.Difficulty,
		MaxCombo:           *b.MaxCombo,
		HitLength:          *b.HitLength,
		Ranked:             *b.Ranked,
		RankedStatusFrozen: *b.RankedStatusFrozen,
		LatestUpdate:       *b.LatestUpdate,
	}
}

type beatmapResponse struct {
	common.ResponseBase
	beatmap
}
type beatmapSetResponse struct {
	common.ResponseBase
	Beatmaps []beatmap `json:"beatmaps"`
}

type beatmapSetStatusData struct {
	BeatmapsetID int `json:"beatmapset_id"`
	BeatmapID    int `json:"beatmap_id"`
	RankedStatus int `json:"ranked_status"`
	Frozen       int `json:"frozen"`
}

// BeatmapSetStatusPOST changes the ranked status of a beatmap, and whether
// the beatmap ranked status is frozen. Or freezed. Freezed best meme 2k16
func BeatmapSetStatusPOST(md common.MethodData) common.CodeMessager {
	var req beatmapSetStatusData
	md.RequestData.Unmarshal(&req)

	var miss []string
	if req.BeatmapsetID == 0 && req.BeatmapID == 0 {
		miss = append(miss, "beatmapset_id or beatmap_id")
	}
	if len(miss) != 0 {
		return ErrMissingField(miss...)
	}

	if req.Frozen != 0 && req.Frozen != 1 {
		return common.SimpleResponse(400, "frozen status must be either 0 or 1")
	}
	if req.RankedStatus > 4 || -1 > req.RankedStatus {
		return common.SimpleResponse(400, "ranked status must be 5 < x < -2")
	}

	param := req.BeatmapsetID
	if req.BeatmapID != 0 {
		err := md.DB.QueryRow("SELECT beatmapset_id FROM beatmaps WHERE beatmap_id = ? LIMIT 1", req.BeatmapID).Scan(&param)
		switch {
		case err == sql.ErrNoRows:
			return common.SimpleResponse(404, "That beatmap could not be found!")
		case err != nil:
			md.Err(err)
			return Err500
		}
	}

	md.DB.Exec(`UPDATE beatmaps 
		SET ranked = ?, ranked_status_freezed = ?
		WHERE beatmapset_id = ?`, req.RankedStatus, req.Frozen, param)

	return getSet(md, param)
}

// BeatmapGET retrieves a beatmap.
func BeatmapGET(md common.MethodData) common.CodeMessager {
	if md.C.Query("s") == "" && md.C.Query("b") == "" {
		return common.SimpleResponse(400, "Must pass either querystring param 'b' or 's'")
	}
	setID := common.Int(md.C.Query("s"))
	if setID != 0 {
		return getSet(md, setID)
	}
	beatmapID := common.Int(md.C.Query("b"))
	if beatmapID != 0 {
		return getBeatmap(md, beatmapID)
	}
	return common.SimpleResponse(400, "Please pass either a valid beatmapset ID or a valid beatmap ID")
}

const baseBeatmapSelect = `
SELECT
	beatmap_id, beatmapset_id, beatmap_md5,
	song_name, ar, od, difficulty, max_combo,
	hit_length, ranked, ranked_status_freezed,
	latest_update
FROM beatmaps
`

func getSet(md common.MethodData, setID int) common.CodeMessager {
	rows, err := md.DB.Query(baseBeatmapSelect+"WHERE beatmapset_id = ?", setID)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r beatmapSetResponse
	for rows.Next() {
		var (
			b               beatmap
			rawLatestUpdate int64
		)
		err = rows.Scan(
			&b.BeatmapID, &b.BeatmapsetID, &b.BeatmapMD5,
			&b.SongName, &b.AR, &b.OD, &b.Difficulty, &b.MaxCombo,
			&b.HitLength, &b.Ranked, &b.RankedStatusFrozen,
			&rawLatestUpdate,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		b.LatestUpdate = time.Unix(rawLatestUpdate, 0)
		r.Beatmaps = append(r.Beatmaps, b)
	}
	r.Code = 200
	return r
}

func getBeatmap(md common.MethodData, beatmapID int) common.CodeMessager {
	var (
		b               beatmap
		rawLatestUpdate int64
	)
	err := md.DB.QueryRow(baseBeatmapSelect+"WHERE beatmap_id = ? LIMIT 1", beatmapID).Scan(
		&b.BeatmapID, &b.BeatmapsetID, &b.BeatmapMD5,
		&b.SongName, &b.AR, &b.OD, &b.Difficulty, &b.MaxCombo,
		&b.HitLength, &b.Ranked, &b.RankedStatusFrozen,
		&rawLatestUpdate,
	)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "That beatmap could not be found!")
	case err != nil:
		md.Err(err)
		return Err500
	}
	b.LatestUpdate = time.Unix(rawLatestUpdate, 0)
	var r beatmapResponse
	r.Code = 200
	r.beatmap = b
	return r
}

type beatmapReduced struct {
	BeatmapID          int    `json:"beatmap_id"`
	BeatmapsetID       int    `json:"beatmapset_id"`
	BeatmapMD5         string `json:"beatmap_md5"`
	Ranked             int    `json:"ranked"`
	RankedStatusFrozen int    `json:"ranked_status_frozen"`
}

type beatmapRankedFrozenFullResponse struct {
	common.ResponseBase
	Beatmaps []beatmapReduced `json:"beatmaps"`
}

// BeatmapRankedFrozenFullGET retrieves all beatmaps with a certain
// ranked_status_freezed
func BeatmapRankedFrozenFullGET(md common.MethodData) common.CodeMessager {
	rows, err := md.DB.Query(`
	SELECT beatmap_id, beatmapset_id, beatmap_md5, ranked, ranked_status_freezed
	FROM beatmaps
	WHERE ranked_status_freezed = '1'
	`)
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r beatmapRankedFrozenFullResponse
	for rows.Next() {
		var b beatmapReduced
		err = rows.Scan(&b.BeatmapID, &b.BeatmapsetID, &b.BeatmapMD5, &b.Ranked, &b.RankedStatusFrozen)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Beatmaps = append(r.Beatmaps, b)
	}
	r.Code = 200
	return r
}
