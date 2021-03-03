/**
 * @Time: 2021/3/2 10:35 上午
 * @Author: varluffy
 * @Description: //TODO
 */

package entity

import "time"

type UserOauths struct {
	ID           int64     `gorm:"column:id" json:"id" form:"id"`
	UserId       int64     `gorm:"column:user_id" json:"user_id" form:"user_id"`
	IdentityType int64     `gorm:"column:identity_type" json:"identity_type" form:"identity_type"`
	Identifier   string    `gorm:"column:identifier" json:"identifier" form:"identifier"`
	UnionId      string    `gorm:"column:union_id" json:"union_id" form:"union_id"`
	Credential   string    `gorm:"column:credential" json:"credential" form:"credential"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at" form:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at" form:"updated_at"`
}

func (u UserOauths) TableName() string {
	return "user_oauths"
}
