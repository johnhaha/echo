package echo

import "errors"

var (
	errNoSubscriber = errors.New("😥 no valid subscriber")
	errNotFound     = errors.New("🥲 not found")
)
