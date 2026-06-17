package models

import "time"

type Member struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64     `gorm:"not null;index" json:"user_id"`
	MembershipType string     `gorm:"type:varchar(50);not null;default:basic" json:"membership_type"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Status         int8       `gorm:"default:1" json:"status"`
	Notes          string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy      uint64     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Member) TableName() string {
	return "members"
}
