package models

import "time"

type Schedule struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NewsID       int       `gorm:"uniqueIndex" json:"news_id"`
	Title        string    `gorm:"type:text" json:"title"`
	Date         time.Time `gorm:"not null" json:"date"`
	ScheduleDate string    `gorm:"column:schedule_date" json:"schedule_date"`
	Col1_1       string    `gorm:"column:1_1" json:"1_1"`
	Col1_2       string    `gorm:"column:1_2" json:"1_2"`
	Col2_1       string    `gorm:"column:2_1" json:"2_1"`
	Col2_2       string    `gorm:"column:2_2" json:"2_2"`
	Col3_1       string    `gorm:"column:3_1" json:"3_1"`
	Col3_2       string    `gorm:"column:3_2" json:"3_2"`
	Col4_1       string    `gorm:"column:4_1" json:"4_1"`
	Col4_2       string    `gorm:"column:4_2" json:"4_2"`
	Col5_1       string    `gorm:"column:5_1" json:"5_1"`
	Col5_2       string    `gorm:"column:5_2" json:"5_2"`
	Col6_1       string    `gorm:"column:6_1" json:"6_1"`
	Col6_2       string    `gorm:"column:6_2" json:"6_2"`
}

func (d *Schedule) FormatDate() string {
	return d.Date.Format("02.01.2006 15:04")
}
