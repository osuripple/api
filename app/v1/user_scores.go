package v1

import "time"

type userScore struct {
	ScoreID    int       `json:"score_id"`
	BeatmapMD5 string    `json:"beatmap_md5"`
	Score      int64     `json:"score"`
	MaxCombo   int       `json:"max_combo"`
	FullCombo  bool      `json:"full_combo"`
	Mods       int       `json:"mods"`
	Count300   int       `json:"count_300"`
	Count100   int       `json:"count_100"`
	Count50    int       `json:"count_50"`
	CountGeki  int       `json:"count_geki"`
	CountKatu  int       `json:"count_katu"`
	CountMiss  int       `json:"count_miss"`
	Time       time.Time `json:"time"`
	PlayMode   int       `json:"play_mode"`
	Accuracy   float64   `json:"accuracy"`
}
