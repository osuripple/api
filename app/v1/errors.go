package v1

import (
	"github.com/osuripple/api/common"
)

// Boilerplate errors
var (
	Err500 = common.Response{
		Code: 0,
		Message: "An error occurred. Try again, perhaps?",
	}
)
