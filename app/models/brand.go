package models

import (
	"gorm.io/gorm"
	"time"
)

type Brand struct {
	Id        uint   `json:"ID"   form:"ID"     gorm:"primary_key"`
	Name      string `json:"name" form:"name"   gorm:"unique;not null"  validate:"required,max=128"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
