package v1

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/osuripple/api/common"
)

type friendData struct {
	userData
	IsMutual bool `json:"is_mutual"`
}

// FriendsGET is the API request handler for GET /friends.
// It retrieves an user's friends, and whether the friendship is mutual or not.
func FriendsGET(md common.MethodData) (r common.Response) {
	var myFrienders []int
	myFriendersRaw, err := md.DB.Query("SELECT user1 FROM users_relationships WHERE user2 = ?", md.User.UserID)
	if err != nil {
		md.C.Error(err)
		r = Err500
		return
	}
	defer myFriendersRaw.Close()
	for myFriendersRaw.Next() {
		var i int
		err := myFriendersRaw.Scan(&i)
		if err != nil {
			md.C.Error(err)
			continue
		}
		myFrienders = append(myFrienders, i)
	}
	if err := myFriendersRaw.Err(); err != nil {
		md.C.Error(err)
	}

	// Yes.
	myFriendsQuery := `
SELECT 
	users.id, users.username, users.register_datetime, users.rank, users.latest_activity,
	
	users_stats.username_aka, users_stats.badges_shown,
	users_stats.country, users_stats.show_country
FROM users_relationships
LEFT JOIN users
ON users_relationships.user2 = users.id
LEFT JOIN users_stats
ON users_relationships.user2=users_stats.id
WHERE users_relationships.user1=?
ORDER BY users_relationships.id`

	results, err := md.DB.Query(myFriendsQuery+common.Paginate(md.C.Query("p"), md.C.Query("l")), md.User.UserID)
	if err != nil {
		md.C.Error(err)
		r = Err500
		return
	}

	var myFriends []friendData

	defer results.Close()
	for results.Next() {
		newFriend := friendPuts(md, results)
		for _, uid := range myFrienders {
			if uid == newFriend.ID {
				newFriend.IsMutual = true
				break
			}
		}
		myFriends = append(myFriends, newFriend)
	}
	if err := results.Err(); err != nil {
		md.C.Error(err)
	}

	r.Code = 200
	r.Data = myFriends
	return
}

func friendPuts(md common.MethodData, row *sql.Rows) (user friendData) {
	var err error

	registeredOn := int64(0)
	latestActivity := int64(0)
	var badges string
	var showcountry bool
	err = row.Scan(&user.ID, &user.Username, &registeredOn, &user.Rank, &latestActivity, &user.UsernameAKA, &badges, &user.Country, &showcountry)
	if err != nil {
		md.C.Error(err)
		return
	}

	user.RegisteredOn = time.Unix(registeredOn, 0)
	user.LatestActivity = time.Unix(latestActivity, 0)

	badgesSl := strings.Split(badges, ",")
	for _, badge := range badgesSl {
		if badge != "" && badge != "0" {
			// We are ignoring errors because who really gives a shit if something's gone wrong on our end in this
			// particular thing, we can just silently ignore this.
			nb, err := strconv.Atoi(badge)
			if err == nil && nb != 0 {
				user.Badges = append(user.Badges, nb)
			}
		}
	}

	// If the user wants to stay anonymous, don't show their country.
	// This can be overriden if we have the ReadConfidential privilege and the user we are accessing is the token owner.
	if !(showcountry || (md.User.Privileges.HasPrivilegeReadConfidential() && user.ID == md.User.UserID)) {
		user.Country = "XX"
	}
	return
}
