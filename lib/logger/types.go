package logger

type logData map[string]interface{}

type keyType string

const (
	dataKey  = keyType("logData")
	levelKey = keyType("slogLevel")
)

const (
	UserIDField     = "user_id"
	RequestIDField  = "request_id"
	InstanceIDField = "instance_id"
)
