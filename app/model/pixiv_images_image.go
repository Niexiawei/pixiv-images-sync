package model

import (
	"gorm.io/gorm"
	"pixivImages/database"
)

type PixivImagesImage struct {
	Id          int    `gorm:"column:id"`
	PixivId     int64  `gorm:"column:pixiv_id"`
	OriginUrl   string `gorm:"column:origin_url"`
	ItemId      string `gorm:"column:item_id"`
	DownloadUrl string `gorm:"column:download_url"`
}

func (PixivImagesImage) TableName() string {
	return "pixiv_images_image"
}

func GetPixivImagesImageDb() *gorm.DB {
	return database.GetMysqlDb().Model(&PixivImagesImage{})
}

func (p *PixivImagesImage) BeforeCreate(tx *gorm.DB) (err error) {

}
