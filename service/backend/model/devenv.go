package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Devenv struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;"`
	UserID    uint      `json:"-"`
	Public    bool      `json:"public"`
	Name      string    `json:"name"`
	BuildCmd  string    `json:"buildCmd"`
	RunCmd    string    `json:"runCmd"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

func (devenv *Devenv) BeforeCreate(tx *gorm.DB) (err error) {
	devenv.ID = uuid.New().String()
	return
}
