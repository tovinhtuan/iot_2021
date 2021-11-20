package model

import "time"

type Sensor struct {
	Id          int64   `gorm:"primaryKey" json:"id"`
	Flight      bool    `gorm:"type:boolean" json:"flight"`
	Temperature float64 `gorm:"type:float" json:"temperature"`
	Humidity    float64 `gorm:"type:float" json:"humidity"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
