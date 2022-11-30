package pixiv_images_save

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"pixivImages/app/model"
	"pixivImages/app/service/onedrive"
	"pixivImages/app/service/pixiv"
	"strconv"
	"strings"
	"time"
)

const (
	PixivImagesSave = "pixiv:images:save"
	BigFile         = 4 * 1024 * 1024
)

var (
	UploadError = errors.New("图片上传失败")
)

type SaveImageTask struct {
	pixiv.Image
	Date time.Time
}

func NewTask(image pixiv.Image, date time.Time) *asynq.Task {
	saveImage := SaveImageTask{
		Date:  date,
		Image: image,
	}
	payload, _ := json.Marshal(saveImage)
	return asynq.NewTask(PixivImagesSave, payload, asynq.MaxRetry(5))
}

type Process struct {
	oneDriveClient *onedrive.OneDrive
	pixivClient    *pixiv.Pixiv
}

func (p *Process) imagesUpload(url pixiv.ImageUrl, IllustId string) (string, error) {
	content, err := p.pixivClient.GetImageContent(url.Url, IllustId)
	if err != nil {
		return "", err
	}

	path := "/" + strings.Join(url.GetPath(), "/")
	result := &onedrive.UploadFileResult{}
	if len(content) > BigFile {
		result, err = p.oneDriveClient.UploadBigFile(url.Name, path, bytes.NewReader(content), int64(len(content)))
	} else {
		result, err = p.oneDriveClient.UploadFile(url.Name, path, bytes.NewReader(content))
	}

	if err != nil {
		return "", errors.WithMessage(UploadError, err.Error())
	}

	return result.Id, nil
}

func (p *Process) imageSave(image SaveImageTask) (*model.PixivImages, error) {
	var saveUploadErr error
	images := model.PixivImages{
		Title:    image.Title,
		IllustId: strconv.Itoa(image.IllustId),
	}
	db1 := model.GetPixivImagesDb()
	if err := db1.Where(map[string]interface{}{
		"illust_id": image.IllustId,
	}).FirstOrCreate(&images).Error; err != nil {
		return nil, err
	}

	for _, url := range image.Urls {
		_image := &model.PixivImagesImage{
			PixivId:     images.PixivId,
			OriginUrl:   url.Url,
			ItemId:      "",
			DownloadUrl: "",
		}
		singleImageDb := model.GetPixivImagesImageDb()
		if err := singleImageDb.First(&_image, "origin_url = ?", url.Url).Error; !(err != nil && errors.Is(err, gorm.ErrRecordNotFound)) {
			continue
		}
		id, err := p.imagesUpload(url, strconv.Itoa(image.IllustId))
		if err != nil {
			saveUploadErr = err
			continue
		}
		_image.ItemId = id
		model.GetPixivImagesImageDb().Create(&_image)
	}
	return &images, saveUploadErr
}

func ProcessTask(ctx context.Context, task *asynq.Task) error {
	var image SaveImageTask
	if err := json.Unmarshal(task.Payload(), &image); err != nil {
		return err
	}
	run := &Process{
		onedrive.NewOneDrive(),
		pixiv.NewPixiv(),
	}
	images, err := run.imageSave(image)

	if err != nil {
		return err
	}

	date := image.Date.Format("2006-01-02")

	rank := model.PixivImagesRank{
		PixivId:  images.PixivId,
		IllustId: images.IllustId,
		Rank:     image.Rank,
		Date:     date,
	}

	model.GetPixivImagesRankDb().Where(map[string]interface{}{
		"rank": image.Rank,
		"date": date,
	}).FirstOrCreate(&rank)
	return nil
}
