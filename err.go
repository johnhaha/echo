package echo

import "errors"

var (
	errInvalidJson  = errors.New("😥 invalid json data")
	errNoSubscriber = errors.New("😥 no valid subscriber")
)
