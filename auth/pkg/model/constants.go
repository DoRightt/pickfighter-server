package model

// Services names used in the application
const (
	AuthService = "auth"
)

// ContextField represents keys used to store values in the context.
type ContextField string

// Constants defining keys for commonly used context values.
const (
	ContextUserId          ContextField = "user_id"
	ContextFlags           ContextField = "flags"
	ContextClaim           ContextField = "root_claim"
	ContextNamespaceClaims ContextField = "ns_claims"
	ContextJWTPointer      ContextField = "jwt_pointer"
)
