package model

import (
	"gorm.io/gorm"
	"pixivImages/database"
	"pixivImages/utils"
	"time"
)

type PixivImages struct {
	Id        int       `gorm:"column:id"`
	PixivId   int64     `gorm:"column:pixiv_id"`
	Title     string    `gorm:"column:title"`
	IllustId  string    `gorm:"column:illust_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PixivImages) TableName() string {
	return "pixiv_images"
}

func GetPixivImagesDb() *gorm.DB {
	return database.GetMysqlDb().Model(&PixivImages{})
}

func (p *PixivImages) BeforeCreate(tx *gorm.DB) (err error) {
	p.PixivId = utils.GetNextId()
	return
}
