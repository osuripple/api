package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/osuripple/api/common"
)

// Handle404 handles requests with no implemented handlers.
func Handle404(c *gin.Context) {
	c.IndentedJSON(404, common.Response{
		Code:    404,
		Message: "Oh dear... that API request could not be found! Perhaps the API is not up-to-date? Either way, have a surprise!",
		Data:    surpriseMe(),
	})
}
