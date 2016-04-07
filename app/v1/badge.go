package v1

import (
	"database/sql"

	"github.com/osuripple/api/common"
)

type badgeData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// BadgeByIDGET is the handler for /badge/:id
func BadgeByIDGET(md common.MethodData) (r common.Response) {
	b := badgeData{}
	err := md.DB.QueryRow("SELECT id, name, icon FROM badges WHERE id=? LIMIT 1", md.C.Param("id")).Scan(&b.ID, &b.Name, &b.Icon)
	switch {
	case err == sql.ErrNoRows:
		r.Code = 404
		r.Message = "No such badge was found"
		return
	case err != nil:
		md.C.Error(err)
		r = Err500
		return
	}
	r.Code = 200
	r.Data = b
	return
}

// BadgesGET retrieves all the badges on this ripple instance.
func BadgesGET(md common.MethodData) (r common.Response) {
	var badges []badgeData
	rows, err := md.DB.Query("SELECT id, name, icon FROM badges")
	if err != nil {
		md.C.Error(err)
		r = Err500
		return
	}
	defer rows.Close()
	for rows.Next() {
		nb := badgeData{}
		err = rows.Scan(&nb.ID, &nb.Name, &nb.Icon)
		if err != nil {
			md.C.Error(err)
		}
		badges = append(badges, nb)
	}
	if err := rows.Err(); err != nil {
		md.C.Error(err)
	}
	r.Code = 200
	r.Data = badges
	return
}
