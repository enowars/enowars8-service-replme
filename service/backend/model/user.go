package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"unique"`
	Password  string         `json:"-"`
	Devenvs   []Devenv       `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created"`
	UpdatedAt time.Time      `json:"updated"`
	DeletedAt gorm.DeletedAt `json:"deleted" gorm:"index"`
}
