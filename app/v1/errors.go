package v1

import (
	"strings"

	"github.com/osuripple/api/common"
)

// Boilerplate errors
var (
	Err500 = common.Response{
		Code:    500,
		Message: "An error occurred. Try again, perhaps?",
	}
	ErrBadJSON = common.Response{
		Code:    400,
		Message: "There was an error processing your JSON data.",
	}
)

// ErrMissingField generates a response to a request when some fields in the JSON are missing.
func ErrMissingField(missingFields ...string) common.Response {
	return common.Response{
		Code:    422, // http://stackoverflow.com/a/10323055/5328069
		Message: "Missing fields: " + strings.Join(missingFields, ", ") + ".",
	}
}
