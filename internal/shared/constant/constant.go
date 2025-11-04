package constant

type contextKey string

const (
	SpanTypeProccess  = "process"
	SpanActionExecute = "execute"
	SpanTypeDb        = "db"
	SpanSubTypeDb     = "postgresql"
	SpanActionQuery   = "query"

	ContextKeyProcessId      contextKey = "process-id"
	ContextKeyProductName    contextKey = "product-name"
	ContextKeyProcessIdStr              = "process-id"
	ContextKeyProductNameStr            = "product-name"

	UserTypeSystemAdmin = "SYSTEM_ADMIN"
	UserTypeProjectUser = "PROJECT_USER"
)
