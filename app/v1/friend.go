package v1

import (
	"database/sql"
	"strconv"
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
	myFriendersRaw, err := md.DB.Query("SELECT user1 FROM users_relationships WHERE user2 = ?", md.ID())
	if err != nil {
		md.Err(err)
		r = Err500
		return
	}
	defer myFriendersRaw.Close()
	for myFriendersRaw.Next() {
		var i int
		err := myFriendersRaw.Scan(&i)
		if err != nil {
			md.Err(err)
			continue
		}
		myFrienders = append(myFrienders, i)
	}
	if err := myFriendersRaw.Err(); err != nil {
		md.Err(err)
	}

	// Yes.
	myFriendsQuery := `
SELECT 
	users.id, users.username, users.register_datetime, users.rank, users.latest_activity,
	
	users_stats.username_aka,
	users_stats.country, users_stats.show_country
FROM users_relationships
LEFT JOIN users
ON users_relationships.user2 = users.id
LEFT JOIN users_stats
ON users_relationships.user2=users_stats.id
WHERE users_relationships.user1=?
ORDER BY users_relationships.id`

	results, err := md.DB.Query(myFriendsQuery+common.Paginate(md.C.Query("p"), md.C.Query("l")), md.ID())
	if err != nil {
		md.Err(err)
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
		md.Err(err)
	}

	r.Code = 200
	r.Data = myFriends
	return
}

func friendPuts(md common.MethodData, row *sql.Rows) (user friendData) {
	var err error

	registeredOn := int64(0)
	latestActivity := int64(0)
	var showcountry bool
	err = row.Scan(&user.ID, &user.Username, &registeredOn, &user.Rank, &latestActivity, &user.UsernameAKA, &user.Country, &showcountry)
	if err != nil {
		md.Err(err)
		return
	}

	user.RegisteredOn = time.Unix(registeredOn, 0)
	user.LatestActivity = time.Unix(latestActivity, 0)

	// If the user wants to stay anonymous, don't show their country.
	// This can be overriden if we have the ReadConfidential privilege and the user we are accessing is the token owner.
	if !(showcountry || (md.User.Privileges.HasPrivilegeReadConfidential() && user.ID == md.ID())) {
		user.Country = "XX"
	}
	return
}

type friendsWithData struct {
	Friends bool `json:"friend"`
	Mutual  bool `json:"mutual"`
}

// FriendsWithGET checks the current user is friends with the one passed in the request path.
func FriendsWithGET(md common.MethodData) (r common.Response) {
	r.Code = 200
	var d friendsWithData
	uid, err := strconv.Atoi(md.C.Param("id"))
	if err != nil {
		r.Data = d
		return
	}
	err = md.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users_relationships WHERE user1 = ? AND user2 = ? LIMIT 1), EXISTS(SELECT 1 FROM users_relationships WHERE user2 = ? AND user1 = ? LIMIT 1)", md.ID(), uid, md.ID(), uid).Scan(&d.Friends, &d.Mutual)
	if err != sql.ErrNoRows && err != nil {
		md.Err(err)
		r = Err500
		return
	}
	r.Data = d
	return
}

// FriendsAddGET is the GET version of FriendsAddPOST.
func FriendsAddGET(md common.MethodData) common.Response {
	uidS := md.C.Param("id")
	uid, err := strconv.Atoi(uidS)
	if err != nil {
		return common.Response{
			Code:    400,
			Message: "Nope. That's not a number.",
		}
	}
	return addFriend(md, uid)
}

type friendAddPOSTData struct {
	UserID int `json:"user_id"`
}

// FriendsAddPOST allows for adding friends. Yup. Easy as that.
func FriendsAddPOST(md common.MethodData) (r common.Response) {
	d := friendAddPOSTData{}
	err := md.RequestData.Unmarshal(&d)
	if err != nil {
		md.Err(err)
		r = Err500
		return
	}
	return addFriend(md, d.UserID)
}

func addFriend(md common.MethodData, u int) (r common.Response) {
	if md.ID() == u {
		r.Code = 400
		r.Message = "Just so you know: you can't add yourself to your friends."
		return
	}
	if !userExists(md, u) {
		r.Code = 404
		r.Message = "I'd also like to be friends with someone who doesn't even exist (???), however that's NOT POSSIBLE."
		return
	}
	var (
		relExists bool
		isMutual  bool
	)
	err := md.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users_relationships WHERE user1 = ? AND user2 = ?), EXISTS(SELECT 1 FROM users_relationships WHERE user2 = ? AND user1 = ?)", md.ID(), u, md.ID(), u).Scan(&relExists, &isMutual)
	if err != nil && err != sql.ErrNoRows {
		md.Err(err)
		r = Err500
		return
	}
	if !relExists {
		_, err := md.DB.Exec("INSERT INTO users_relationships(user1, user2) VALUES (?, ?)", md.User.UserID, u)
		if err != nil {
			md.Err(err)
			r = Err500
			return
		}
	}
	r.Code = 200
	r.Data = friendsWithData{
		Friends: true,
		Mutual:  isMutual,
	}
	return
}

// userExists makes sure an user exists.
func userExists(md common.MethodData, u int) (r bool) {
	err := md.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", u).Scan(&r)
	if err != nil && err != sql.ErrNoRows {
		md.Err(err)
	}
	return
}

// FriendsDelGET is the GET version of FriendDelPOST.
func FriendsDelGET(md common.MethodData) common.Response {
	uidS := md.C.Param("id")
	uid, err := strconv.Atoi(uidS)
	if err != nil {
		return common.Response{
			Code:    400,
			Message: "Nope. That's not a number.",
		}
	}
	return delFriend(md, uid)
}

// FriendsDelPOST allows for deleting friends.
func FriendsDelPOST(md common.MethodData) (r common.Response) {
	d := friendAddPOSTData{}
	err := md.RequestData.Unmarshal(&d)
	if err != nil {
		md.Err(err)
		r = Err500
		return
	}
	return delFriend(md, d.UserID)
}

func delFriend(md common.MethodData, u int) common.Response {
	_, err := md.DB.Exec("DELETE FROM users_relationships WHERE user1 = ? AND user2 = ?", md.ID(), u)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return common.Response{
		Code: 200,
		Data: friendsWithData{
			Friends: false,
			Mutual:  false,
		},
	}
}
