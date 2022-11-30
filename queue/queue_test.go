package queue

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pixivImages/app/model"
	pixiv2 "pixivImages/app/service/pixiv"
	"pixivImages/config"
	"pixivImages/database"
	"pixivImages/queue/tasks/pixiv_images_save"
	"testing"
	"time"
)

func Test_66(t *testing.T) {
	config.LoadConfig()
	database.InitMysql()
	_image := &model.PixivImagesImage{
		PixivId:     1597096997435740160,
		OriginUrl:   "https://i.pximg.net/img-original/img/2022/11/14/08/09/40/102781051_p0.png",
		ItemId:      "",
		DownloadUrl: "",
	}
	singleImageDb := model.GetPixivImagesImageDb()
	if err := singleImageDb.First(&_image, "origin_url = ?", _image.OriginUrl).Error; !(err != nil && errors.Is(err, gorm.ErrRecordNotFound)) {
		fmt.Println("存在")
		return
	}
	fmt.Println("不存在")
}

func Test_Inspector(t *testing.T) {
	config.LoadConfig()
	InitClient()
	info, err := GetInspector().GetTaskInfo("default", "d346fbd3-68cf-4895-a78c-9321873946fe")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(info.Type)
}

func Test_RunArchivedTask(t *testing.T) {
	config.LoadConfig()
	InitClient()
	err := GetInspector().RunTask("default", "0d46e2ab-0bb2-43a1-a5dc-1e534897936a")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func TestQueue(t *testing.T) {
	config.LoadConfig()
	InitClient()

	pixiv := pixiv2.NewPixiv()
	date := time.Date(2022, 11, 24, 0, 0, 0, 0, time.Local)
	images, err := pixiv.RankImageUrls(pixiv2.WithRankListDate(date), pixiv2.WithRankQueryMaxQuantity(100))
	if err != nil {
		t.Error(err)
		return
	}
	taskClient := GetClient()
	for _, image := range images {
		info, err := taskClient.Enqueue(pixiv_images_save.NewTask(image, date))
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Printf("%+v\n", info.ID)
	}
}
