// middleware/context.go
package middleware

type ContextKey string

const UserIdContextKey ContextKey = "userID"

const SessionIdContextKey ContextKey = "authStatus"
