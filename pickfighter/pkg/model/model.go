package model

type HealthStatus struct {
	AppDevVersion string `json:"app_dev_version"`
	AppName       string `json:"app_name"`
	AppRunDate    int64  `json:"app_run_date"`
	AppTimeAlive  int64  `json:"app_time_alive"`
	Healthy       bool   `json:"healthy"`
	Message       string `json:"message"`
	Timestamp     string `json:"timestamp"`
}
