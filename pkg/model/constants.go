package model

// Services names used in the application
const (
	GatewayService  = "gateway"
	AuthService     = "auth"
	EventService    = "events"
	FightersService = "fighters"
)

// ContextField represents keys used to store values in the context.
type ContextField string

// Constants defining keys for commonly used context values.
const (
	ContextUserId          ContextField = "user_id"
	ContextFlags           ContextField = "flags"
	ContextClaim           ContextField = "root_claim"
	ContextJWTPointer      ContextField = "jwt_pointer"
)

const (
	DefaultPasswordLength = 64
	DefaultSaltLength     = 32
)
