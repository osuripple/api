package websockets

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/thehowl/go-osuapi.v1"
	"zxq.co/ripple/rippleapi/app/v1"
	"zxq.co/x/getrank"
)

type subscribeScoresUser struct {
	User  int   `json:"user"`
	Modes []int `json:"modes"`
}

// SubscribeScores subscribes a connection to score updates.
func SubscribeScores(c *conn, message incomingMessage) {
	var ssu []subscribeScoresUser
	err := json.Unmarshal(message.Data, &ssu)
	if err != nil {
		c.WriteJSON(TypeInvalidMessage, err.Error())
		return
	}

	scoreSubscriptionsMtx.Lock()

	var found bool
	for idx, el := range scoreSubscriptions {
		// already exists, change the users
		if el.Conn.ID == c.ID {
			found = true
			scoreSubscriptions[idx].Users = ssu
		}
	}

	// if it was not found, we need to add it
	if !found {
		scoreSubscriptions = append(scoreSubscriptions, scoreSubscription{c, ssu})
	}

	scoreSubscriptionsMtx.Unlock()

	c.WriteJSON(TypeSubscribedToScores, ssu)
}

type scoreSubscription struct {
	Conn  *conn
	Users []subscribeScoresUser
}

var scoreSubscriptions []scoreSubscription
var scoreSubscriptionsMtx = new(sync.RWMutex)

func scoreRetriever() {
	ps, err := red.Subscribe("api:score_submission")
	if err != nil {
		fmt.Println(err)
	}
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go handleNewScore(msg.Payload)
	}
}

type score struct {
	v1.Score
	UserID int `json:"user_id"`
}

func handleNewScore(id string) {
	defer catchPanic()
	var s score
	err := db.Get(&s, `
SELECT
	id, beatmap_md5, score, max_combo, full_combo, mods,
	300_count, 100_count, 50_count, gekis_count, katus_count, misses_count,
	time, play_mode, accuracy, pp, completed, userid AS user_id
FROM scores WHERE id = ?`, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Rank = strings.ToUpper(getrank.GetRank(
		osuapi.Mode(s.PlayMode),
		osuapi.Mods(s.Mods),
		s.Accuracy,
		s.Count300,
		s.Count100,
		s.Count50,
		s.CountMiss,
	))
	scoreSubscriptionsMtx.RLock()
	cp := make([]scoreSubscription, len(scoreSubscriptions))
	copy(cp, scoreSubscriptions)
	scoreSubscriptionsMtx.RUnlock()

	for _, el := range cp {
		if len(el.Users) > 0 && !scoreUserValid(el.Users, s) {
			continue
		}

		el.Conn.WriteJSON(TypeNewScore, s)
	}
}

func scoreUserValid(users []subscribeScoresUser, s score) bool {
	for _, u := range users {
		if u.User == s.UserID {
			if len(u.Modes) > 0 {
				if !inModes(u.Modes, s.PlayMode) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func inModes(modes []int, i int) bool {
	for _, m := range modes {
		if m == i {
			return true
		}
	}
	return false
}

func catchPanic() {
	r := recover()
	if r != nil {
		fmt.Println(r)
		// TODO: sentry
	}
}
