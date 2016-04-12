package v1

import "github.com/osuripple/api/common"

type setAllowedData struct {
	UserID  int `json:"user_id"`
	Allowed int `json:"allowed"`
}

// UserManageSetAllowedPOST allows to set the allowed status of an user.
func UserManageSetAllowedPOST(md common.MethodData) (r common.Response) {
	data := setAllowedData{}
	if err := md.RequestData.Unmarshal(&data); err != nil {
		r = ErrBadJSON
		return
	}
	if data.Allowed < 0 || data.Allowed > 2 {
		r.Code = 400
		r.Message = "Allowed status must be between 0 and 2"
		return
	}
	_, err := md.DB.Exec("UPDATE users SET allowed = ? WHERE id = ?", data.Allowed, data.UserID)
	if err != nil {
		md.Err(err)
		r = Err500
		return
	}
	query := `
SELECT users.id, users.username, register_datetime, rank,
	latest_activity, users_stats.username_aka,
	users_stats.country, users_stats.show_country
FROM users
LEFT JOIN users_stats
ON users.id=users_stats.id
WHERE users.id=?
LIMIT 1`
	r = userPuts(md, md.DB.QueryRow(query, data.UserID))
	return
}
