package git

import "errors"

var (
	NO_DIFFERENCES_ERROR error = errors.New("branches have no differences")
)
