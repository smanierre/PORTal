package serverutils

type ContextKey string

const (
	MemberContextKey  ContextKey = "member"
	SessionCookieName            = "session_id"
)
