package v1

import (
	"git.zxq.co/ripple/rippleapi/common"
	"github.com/gin-gonic/gin"
)

type response404 struct {
	common.ResponseBase
	Cats string `json:"cats"`
}

// Handle404 handles requests with no implemented handlers.
func Handle404(c *gin.Context) {
	c.IndentedJSON(404, response404{
		ResponseBase: common.ResponseBase{
			Code:    404,
			Message: "Oh dear... that API request could not be found! Perhaps the API is not up-to-date? Either way, have a surprise!",
		},
		Cats: surpriseMe(),
	})
}
