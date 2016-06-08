package v1

import "git.zxq.co/ripple/rippleapi/common"

type beatmap struct {
	BeatmapID          int     `json:"beatmap_id"`
	BeatmapsetID       int     `json:"beatmapset_id"`
	BeatmapMD5         string  `json:"beatmap_md5"`
	SongName           string  `json:"song_name"`
	AR                 float32 `json:"ar"`
	OD                 float32 `json:"od"`
	Difficulty         float64 `json:"difficulty"`
	MaxCombo           int     `json:"max_combo"`
	HitLength          int     `json:"hit_length"`
	Ranked             int     `json:"ranked"`
	RankedStatusFrozen int     `json:"ranked_status_frozen"`
	LatestUpdate       int     `json:"latest_update"`
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
	LatestUpdate       *int
}

func (b *beatmapMayOrMayNotExist) toBeatmap() *beatmap {
	if b.BeatmapID == nil {
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
	BeatmapSetID int `json:"beatmapset_id"`
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
	if req.BeatmapSetID == 0 && req.BeatmapID == 0 {
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

	var (
		whereClause = "beatmapset_id"
		param       = req.BeatmapSetID
	)
	if req.BeatmapID != 0 {
		whereClause = "beatmap_id"
		param = req.BeatmapID
	}

	md.DB.Exec(`UPDATE beatmaps 
		SET ranked = ?, ranked_status_freezed = ?
		WHERE `+whereClause+` = ?`, req.RankedStatus, req.Frozen, param)

	// TODO: replace with beatmapSetResponse when implemented
	return common.ResponseBase{
		Code: 200,
	}
}
