package service_images

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pixivImages/app/model"
	"pixivImages/app/service/onedrive"
	"pixivImages/database"
	"strings"
	"sync"
)

func RefreshPixivDownloadUrl() {
	id := 0
	oneDriveClient := onedrive.NewOneDrive()
	refresh := func(result []model.PixivImagesImage) error {
		wg := sync.WaitGroup{}
		wg.Add(len(result))
		var (
			files [][]interface{}
			idS   []int
		)
		for _, val := range result {
			go func(image model.PixivImagesImage) {
				url, err := oneDriveClient.GetDownloadUrl(image.ItemId)
				if err == nil {
					files = append(files, []interface{}{image.Id, url})
				}
				wg.Done()
			}(val)
		}
		wg.Wait()

		filesStr := ""
		for _, v := range files {
			filesStr = filesStr + fmt.Sprintf("WHEN %d THEN '%s' \n", v[0].(int), v[1].(string))
		}

		idS = slice.Map(files, func(index int, item []interface{}) int {
			return item[0].(int)
		})

		inStr := ""
		slice.ForEach(idS, func(index int, item int) {
			inStr = inStr + fmt.Sprintf("%d,", item)
		})
		inStr = strings.TrimRight(inStr, ",")
		inStr = "(" + inStr + ")"
		sql := fmt.Sprintf("UPDATE pixiv_images_image SET download_url = CASE id \n %s END WHERE id IN %s",
			filesStr,
			inStr,
		)
		return database.GetMysqlDb().Exec(sql).Error
	}

	var result []model.PixivImagesImage
	for {
		dbResult := model.GetPixivImagesImageDb().Where("id > ?", id).Order("id asc").Limit(30).Scan(&result)
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			break
		}
		if dbResult.RowsAffected <= 0 {
			break
		}
		_ = refresh(result)
		id = result[len(result)-1].Id
	}
}
