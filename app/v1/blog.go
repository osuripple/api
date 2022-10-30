package v1

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"zxq.co/ripple/rippleapi/common"
)

// there's gotta be a better way

type blogPost struct {
	ID          string    `json:"id"`
	Creator     blogUser  `json:"creator"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ImportedURL string    `json:"imported_url"`
	UniqueSlug  string    `json:"unique_slug"`

	Snippet     string  `json:"snippet"`
	WordCount   int     `json:"word_count"`
	ReadingTime float64 `json:"reading_time"`
}

type blogUser struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type blogPostsResponse struct {
	common.ResponseBase
	Posts []blogPost `json:"posts"`
}

// consts for the medium API
const (
	mediumFeed              = `https://blog.ripple.moe/feed`
	mediumAPIResponsePrefix = `])}while(1);</x>`
	mediumAPIAllPosts       = `https://blog.ripple.moe/latest?format=json`
)

func init() {
	gob.Register([]blogPost{})
}

// BlogPostsGET retrieves the latest blog posts on the Ripple blog.
func BlogPostsGET(md common.MethodData) common.CodeMessager {
	// check if posts are cached in redis
	res := md.R.Get("api:blog_posts").Val()
	if res != "" {
		// decode values
		posts := make([]blogPost, 0, 20)
		err := gob.NewDecoder(strings.NewReader(res)).Decode(&posts)
		if err != nil {
			md.Err(err)
			return Err500
		}

		// create response and return
		var r blogPostsResponse
		r.Code = 200
		r.Posts = blogLimit(posts, md.Query("l"))
		return r
	}

	// get data from medium rss feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(mediumFeed)
	if err != nil {
		md.Err(fmt.Errorf("rss feed parse url: %w", err))
		return Err500
	}

	// create posts slice and fill it up with converted posts from the medium
	// API
	posts := make([]blogPost, len(feed.Items))
	for idx, mp := range feed.Items {
		var p blogPost

		// convert structs
		p.ID = mp.GUID
		p.Title = mp.Title
		if mp.PublishedParsed != nil {
			p.CreatedAt = *mp.PublishedParsed
		}
		if mp.UpdatedParsed != nil {
			p.UpdatedAt = *mp.UpdatedParsed
		}
		p.ImportedURL = mp.Link
		// p.UniqueSlug = mp.UniqueSlug

		// cr := mResp.Payload.References.User[mp.CreatorID]
		// p.Creator.UserID = cr.UserID
		// p.Creator.Name = cr.Name
		// p.Creator.Username = cr.Username
		if len(mp.Authors) > 0 {
			p.Creator.Name = mp.Authors[0].Name
		}

		// p.Snippet = mp.Virtuals.Subtitle
		// p.WordCount = mp.Virtuals.WordCount
		plainContent := bluemonday.StripTagsPolicy().Sanitize(mp.Content)
		p.Snippet = plainContent[:200]
		if len(p.Snippet) >= 200 {
			p.Snippet += "..."
		}
		p.WordCount = len(plainContent)
		// p.ReadingTime = mp.Virtuals.ReadingTime

		posts[idx] = p
	}

	// save in redis
	bb := new(bytes.Buffer)
	err = gob.NewEncoder(bb).Encode(posts)
	if err != nil {
		md.Err(err)
		return Err500
	}
	md.R.Set("api:blog_posts", bb.Bytes(), time.Minute*5)

	var r blogPostsResponse
	r.Code = 200
	r.Posts = blogLimit(posts, md.Query("l"))
	return r
}

func blogLimit(posts []blogPost, s string) []blogPost {
	i := common.Int(s)
	if i >= len(posts) || i < 1 {
		return posts
	}
	return posts[:i]
}
