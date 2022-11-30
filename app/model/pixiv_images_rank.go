package model

import (
	"gorm.io/gorm"
	"pixivImages/database"
	"time"
)

type PixivImagesRank struct {
	Id        int       `gorm:"column:id"`
	PixivId   int64     `gorm:"column:pixiv_id"`
	IllustId  string    `gorm:"column:illust_id"`
	Rank      int       `gorm:"column:rank"`
	Date      string    `gorm:"column:date"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PixivImagesRank) TableName() string {
	return "pixiv_images_rank"
}

func GetPixivImagesRankDb() *gorm.DB {
	return database.GetMysqlDb().Model(&PixivImagesRank{})
}
