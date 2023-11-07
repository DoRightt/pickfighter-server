package model

type ListRequest struct {
	Limit        int32 `json:"limit" yaml:"limit"`
	Offset       int32 `json:"offset" yaml:"offset"`
	CreatedFrom  int64 `json:"created_at" yaml:"created_at"`
	CreatedUntil int64 `json:"created_until" yaml:"created_until"`
	UpdatedFrom  int64 `json:"updated_at" yaml:"updated_at"`
	UpdatedUntil int64 `json:"updated_until" yaml:"updated_until"`
}