package v1

import (
	"database/sql"

	"github.com/osuripple/api/common"
)

type singleBadge struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type badgeData struct {
	common.ResponseBase
	singleBadge
}

// BadgeByIDGET is the handler for /badge/:id
func BadgeByIDGET(md common.MethodData) common.CodeMessager {
	var b badgeData
	err := md.DB.QueryRow("SELECT id, name, icon FROM badges WHERE id=? LIMIT 1", md.C.Param("id")).Scan(&b.ID, &b.Name, &b.Icon)
	switch {
	case err == sql.ErrNoRows:
		return common.SimpleResponse(404, "No such badge was found")
	case err != nil:
		md.Err(err)
		return Err500
	}
	b.Code = 200
	return b
}

type multiBadgeData struct {
	common.ResponseBase
	Badges []singleBadge `json:"badges"`
}

// BadgesGET retrieves all the badges on this ripple instance.
func BadgesGET(md common.MethodData) common.CodeMessager {
	var r multiBadgeData
	rows, err := md.DB.Query("SELECT id, name, icon FROM badges")
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
