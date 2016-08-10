package v1

import (
	"time"

	"git.zxq.co/ripple/rippleapi/common"
)

type blogPost struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
	Author  userData  `json:"author"`
}

type blogPostsResponse struct {
	common.ResponseBase
	Posts []blogPost `json:"posts"`
}

// BlogPostsGET retrieves the latest blog posts on the Ripple blog.
func BlogPostsGET(md common.MethodData) common.CodeMessager {
	var and string
	var params []interface{}
	if md.C.Query("id") != "" {
		and = "b.id = ?"
		params = append(params, md.C.Query("id"))
	}
	rows, err := md.DB.Query(`
	SELECT 
		b.id, b.title, b.slug, b.created,
		
		u.id, u.username, s.username_aka, u.register_datetime,
		u.privileges, u.latest_activity, s.country
	FROM anchor_posts b
	LEFT JOIN users u ON b.author = u.id
	LEFT JOIN users_stats s ON b.author = s.id
	WHERE status = "published" `+and+`
	ORDER BY b.id DESC `+common.Paginate(md.C.Query("p"), md.C.Query("l"), 50), params...)
	if err != nil {
		md.Err(err)
		return Err500
	}

	var r blogPostsResponse
	for rows.Next() {
		var post blogPost
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Created,

			&post.Author.ID, &post.Author.Username, &post.Author.UsernameAKA, &post.Author.RegisteredOn,
			&post.Author.Privileges, &post.Author.LatestActivity, &post.Author.Country,
		)
		if err != nil {
			md.Err(err)
			continue
		}
		r.Posts = append(r.Posts, post)
	}
	r.Code = 200

	return r
}

type blogPostContent struct {
	common.ResponseBase
	Content string `json:"content"`
}

// BlogPostsContentGET retrieves the content of a specific blog post.
func BlogPostsContentGET(md common.MethodData) common.CodeMessager {
	field := "markdown"
	if _, present := md.C.GetQuery("html"); present {
		field = "html"
	}
	var (
		by  string
		val string
	)
	switch {
	case md.C.Query("slug") != "":
		by = "slug"
		val = md.C.Query("slug")
	case md.C.Query("id") != "":
		by = "id"
		val = md.C.Query("id")
	default:
		return ErrMissingField("id|slug")
	}
	var r blogPostContent
	err := md.DB.QueryRow("SELECT "+field+" FROM anchor_posts WHERE "+by+" = ? AND status = 'published'", val).Scan(&r.Content)
	if err != nil {
		return common.SimpleResponse(404, "no blog post found")
	}
	r.Code = 200
	return r
}
