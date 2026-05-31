package types

type CreateTaskRequest struct {
	Type       string                 `json:"type" binding:"required"`
	Payload    map[string]interface{} `json:"payload" binding:"required"`
	Priority   int                    `json:"priority"`
	MaxRetries int                    `json:"max_retries"`
}
