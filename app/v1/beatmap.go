package v1

import "git.zxq.co/ripple/rippleapi/common"

type beatmap struct {
	BeatmapID          int     `json:"beatmap_id"`
	BeatmapsetID       int     `json:"beatmapset_id"`
	BeatmapMD5         int     `json:"beatmap_md5"`
	SongName           int     `json:"song_name"`
	AR                 float32 `json:"ar"`
	OD                 float32 `json:"od"`
	Difficulty         float64 `json:"difficulty"`
	MaxCombo           int     `json:"max_combo"`
	HitLength          int     `json:"hit_length"`
	BPM                float64 `json:"bpm"`
	Ranked             int     `json:"ranked"`
	RankedStatusFrozen int     `json:"ranked_status_frozen"`
	LatestUpdate       int     `json:"latest_update"`
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
	RankedStatus int `json:"ranked_status"`
	Frozen       int `json:"frozen"`
}

// BeatmapSetStatusPOST changes the ranked status of a beatmap, and whether the beatmap ranked status is frozen. Or freezed. Freezed best meme 2k16
func BeatmapSetStatusPOST(md common.MethodData) common.CodeMessager {
	var req beatmapSetStatusData
	md.RequestData.Unmarshal(&req)

	var miss []string
	if req.BeatmapSetID == 0 {
		miss = append(miss, "beatmapset_id")
	}
	if len(miss) != 0 {
		return ErrMissingField(miss...)
	}

	if req.Frozen != 0 && req.Frozen != 1 {
		return common.SimpleResponse(400, "frozen status must be either 0 or 1")
	}
	if req.RankedStatus > 3 || -2 > req.RankedStatus {
		return common.SimpleResponse(400, "ranked status must be 4 < x < -3")
	}

	md.DB.Exec("UPDATE beatmaps SET ranked = ?, ranked_status_freezed = ? WHERE beatmapset_id = ?", req.RankedStatus, req.Frozen, req.BeatmapSetID)

	// TODO: replace with beatmapSetResponse when implemented
	return common.ResponseBase{
		Code: 200,
	}
}
