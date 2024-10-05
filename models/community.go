package models

import "time"

// Community 社区结构体
type Community struct {
	ID   int64  `json:"community_id" db:"community_id"`
	Name string `json:"community_name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"community_id" db:"community_id"`
	Name         string    `json:"community_name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
