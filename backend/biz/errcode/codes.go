package errcode

const (
	Success = 0

	// 1xxx: auth
	ErrInvalidCredentials = 1001
	ErrTokenExpired       = 1002
	ErrTokenInvalid       = 1003
	ErrUnauthorized       = 1004
	ErrForbidden          = 1005

	// 2xxx: validation
	ErrBadRequest   = 2001
	ErrNotFound     = 2002
	ErrConflict     = 2003
	ErrTooMany      = 2004

	// 3xxx: rate limit & security
	ErrRateLimited = 3001
	ErrBlocked     = 3002

	// 5xxx: internal
	ErrInternal = 5000
	ErrDatabase = 5001
	ErrRedis    = 5002
)
