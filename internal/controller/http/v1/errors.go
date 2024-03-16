package v1

import "fmt"

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
	ErrCannotParseToken  = fmt.Errorf("cannot parse token")
)
