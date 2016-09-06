package beatmapget

import (
	"time"

	"git.zxq.co/ripple/rippleapi/common"
	"gopkg.in/thehowl/go-osuapi.v1"
)

// Set checks if an update is required for all beatmaps in a set.
func Set(s int) error {
	var (
		lastUpdated common.UnixTimestamp
		ranked      int
	)
	err := DB.QueryRow("SELECT last_updated, ranked FROM beatmaps WHERE beatmapset_id = ? LIMIT 1", s).
		Scan(&lastUpdated, &ranked)
	if err != nil {
		return err
	}
	return set(s, lastUpdated, ranked)
}

// Beatmap check if an update is required for all beatmaps in the set
// containing this beatmap.
func Beatmap(b int) error {
	var (
		setID       int
		lastUpdated common.UnixTimestamp
		ranked      int
	)
	err := DB.QueryRow("SELECT beatmapset_id, last_updated, ranked FROM beatmaps WHERE beatmap_id = ? LIMIT 1", b).
		Scan(&setID, &lastUpdated, &ranked)
	if err != nil {
		return err
	}
	return set(setID, lastUpdated, ranked)
}

func set(s int, updated common.UnixTimestamp, ranked int) error {
	expire := Expire
	if ranked == 2 {
		expire *= 6
	}
	if time.Now().Before(time.Time(updated).Add(expire)) {
		return nil
	}
	beatmaps, err := Client.GetBeatmaps(osuapi.GetBeatmapsOpts{
		BeatmapSetID: s,
	})
	if err != nil {
		return err
	}
	for _, beatmap := range beatmaps {
		err := UpdateIfRequired(BeatmapDefiningQuality{
			ID: beatmap.BeatmapID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
