package models

import "time"

type UserModel struct {
	Id         uint64    `gorm:"column:id;primaryKey" json:"id"`
	Username   string    `gorm:"column:username;index" json:"username"`
	Password   string    `gorm:"column:password" json:"password"`
	AvatarUrl  string    `gorm:"column:avatarUrl" json:"avatarUrl"`
	Email      string    `gorm:"column:email;index" json:"email"`
	Phone      string    `gorm:"column:phone;index" json:"phone"`
	UserRole   uint8     `gorm:"column:userRole" json:"userRole"`
	IsDelete   uint8     `gorm:"column:isDelete" json:"isDelete"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:updateTime" json:"updateTime"`
}
