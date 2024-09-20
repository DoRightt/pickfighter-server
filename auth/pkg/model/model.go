package model

// ListRequest provides parameters for listing entities with optional time range filtering.
type ListRequest struct {
	Limit        int32 `json:"limit" yaml:"limit"`
	Offset       int32 `json:"offset" yaml:"offset"`
	CreatedFrom  int64 `json:"created_at" yaml:"created_at"`
	CreatedUntil int64 `json:"created_until" yaml:"created_until"`
	UpdatedFrom  int64 `json:"updated_at" yaml:"updated_at"`
	UpdatedUntil int64 `json:"updated_until" yaml:"updated_until"`
}

type HealthStatus struct {
	AppDevVersion string `json:"app_dev_version"`
	AppName       string `json:"app_name"`
	AppRunDate    int64  `json:"app_run_date"`
	AppTimeAlive  int64  `json:"app_time_alive"`
	Healthy       bool   `json:"healthy"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
}