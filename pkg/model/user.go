package model

type User struct {
	UserId    int32  `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Email     string `json:"email,omitempty" yaml:"email,omitempty"`
	Rank      string `json:"rank,omitempty" yaml:"rank,omitempty"`
	Claim     string `json:"claim,omitempty" yaml:"claim,omitempty"`
	Roles     uint64 `json:"roles,omitempty" yaml:"roles,omitempty"`
	Flags     uint64 `json:"flags,omitempty" yaml:"flags,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty" yaml:"thumbnail,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}
