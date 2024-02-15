package model

const (
	AuthService   = "auth"
	CommonService = "common"
)

type ContextField string

const (
	ContextUserId          ContextField = "user_id"
	ContextFlags           ContextField = "flags"
	ContextClaim           ContextField = "root_claim"
	ContextNamespaceClaims ContextField = "ns_claims"
	ContextJWTPointer      ContextField = "jwt_pointer"
)
