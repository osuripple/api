package common

import (
	"fmt"
	"strconv"
)

// Paginate creates an additional SQL LIMIT clause for paginating.
func Paginate(page, limit string, maxLimit int) string {
	var (
		pInt int
		lInt int
		err  error
	)
	if page == "" {
		pInt = 1
	} else {
		pInt, err = strconv.Atoi(page)
		if err != nil {
			pInt = 1
		}
	}
	if limit == "" {
		lInt = 50
	} else {
		lInt, err = strconv.Atoi(limit)
		if err != nil {
			lInt = 50
		}
	}
	if pInt < 1 {
		pInt = 1
	}
	if lInt < 1 {
		lInt = 50
	}
	if lInt > maxLimit {
		lInt = maxLimit
	}
	start := (pInt - 1) * lInt
	return fmt.Sprintf(" LIMIT %d,%d ", start, lInt)
}
