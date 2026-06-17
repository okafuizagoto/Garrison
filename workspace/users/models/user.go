package models

import "time"

type User struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:varchar(100);not null" json:"name"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Role      string     `gorm:"type:varchar(50);default:member" json:"role"`
	Status    int8       `gorm:"default:1" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
