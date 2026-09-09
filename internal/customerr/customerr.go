package customerr

import "errors"

var ErrNotFoundOrForbidden = errors.New("not found or forbidden")
var ErrTitleEmpty = errors.New("title is empty")
var ErrMaxLimitExceeded = errors.New("limit is exceeded")
