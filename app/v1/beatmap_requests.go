package v1

import (
	"strconv"
	"time"

	"git.zxq.co/ripple/rippleapi/common"
)

type rankRequestsStatusResponse struct {
	common.ResponseBase
	QueueSize       int        `json:"queue_size"`
	MaxPeruser      int        `json:"max_per_user"`
	Submitted       int        `json:"submitted"`
	SubmittedByUser int        `json:"submitted_by_user"`
	CanSubmit       bool       `json:"can_submit"`
	NextExpiration  *time.Time `json:"next_expiration"`
}

// BeatmapRankRequestsStatusGET gets the current status for beatmap ranking requests.
func BeatmapRankRequestsStatusGET(md common.MethodData) common.CodeMessager {
	c := common.GetConf()
	rows, err := md.DB.Query("SELECT userid, time FROM rank_requests WHERE time > ? ORDER BY id ASC LIMIT "+strconv.Itoa(c.RankQueueSize), time.Now().Add(-time.Hour*24).Unix())
	if err != nil {
		md.Err(err)
		return Err500
	}
	var r rankRequestsStatusResponse
	isFirst := true
	for rows.Next() {
		var (
			user      int
			timestamp common.UnixTimestamp
		)
		err := rows.Scan(&user, &timestamp)
		if err != nil {
			md.Err(err)
			continue
		}
		if user == md.ID() {
			r.SubmittedByUser++
		}
		if isFirst {
			x := time.Time(timestamp)
			r.NextExpiration = &x
			isFirst = false
		}
		r.Submitted++
	}
	r.QueueSize = c.RankQueueSize
	r.MaxPeruser = c.BeatmapRequestsPerUser
	r.CanSubmit = r.Submitted < r.QueueSize && r.SubmittedByUser < r.MaxPeruser
	r.Code = 200
	return r
}
