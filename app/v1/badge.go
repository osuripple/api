package v1

import (
	"database/sql"

	"git.zxq.co/ripple/rippleapi/common"
)

type singleBadge struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type multiBadgeData struct {
	common.ResponseBase
	Badges []singleBadge `json:"badges"`
}

// BadgesGET retrieves all the badges on this ripple instance.
func BadgesGET(md common.MethodData) common.CodeMessager {
	var (
		r    multiBadgeData
		rows *sql.Rows
		err  error
	)
	if md.C.Query("id") != "" {
		// TODO(howl): ID validation
		rows, err = md.DB.Query("SELECT id, name, icon FROM badges WHERE id = ?", md.C.Query("id"))
	} else {
		rows, err = md.DB.Query("SELECT id, name, icon FROM badges")
	}
	if err != nil {
		md.Err(err)
		return Err500
	}
	defer rows.Close()
	for rows.Next() {
		nb := singleBadge{}
		err = rows.Scan(&nb.ID, &nb.Name, &nb.Icon)
		if err != nil {
			md.Err(err)
		}
		r.Badges = append(r.Badges, nb)
	}
	if err := rows.Err(); err != nil {
		md.Err(err)
	}
	r.ResponseBase.Code = 200
	return r
}
