package model

import (
	"time"

	"gorm.io/gorm"
)

// 书架
type Shelf struct {
	ID int64 `gorm:"primaryKey"`
	Theme string
	Size int64
	CreateAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// 书
type Book struct {
	ID int64 `gorm:"primaryKey"`
	Author string
	Title string
	CreateAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}