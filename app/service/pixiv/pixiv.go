package pixiv

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"pixivImages/config"
	"pixivImages/logger"
	"pixivImages/utils/gpool"
	"sort"
	"strings"
	"time"
)

var (
	OriginUrlNoFound   = errors.New("未找到原图url")
	RankNotFound       = errors.New("今天没有统计数据")
	ImageDownloadError = errors.New("图片下载失败")
)

type Pixiv struct {
	httpClient   *req.Client
	socks5Config config.Socks5
}

func NewPixiv() *Pixiv {
	socks5Config := config.Get().Socks5
	client := req.C().SetTimeout(30 * time.Second).
		SetProxyURL(fmt.Sprintf("socks5://%s:%s@%s", socks5Config.Username, socks5Config.Password, socks5Config.Addr)).
		SetCommonHeaders(map[string]string{
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36 Edg/104.0.1293.63",
			"x-requested-with": "XMLHttpRequest",
		})
	return &Pixiv{
		httpClient:   client,
		socks5Config: socks5Config,
	}
}

func (client *Pixiv) GetOriginUrl(illustId int) ([]string, error) {
	result := OriginUrlResult{}
	http := client.httpClient.R()
	http.SetHeader("referer", "https://www.pixiv.net/ranking.php?mode=daily&content=illust&lang=zh").
		SetRetryCount(5).
		SetRetryFixedInterval(500 * time.Millisecond)
	_, err := http.SetResult(&result).Get(fmt.Sprintf("https://www.pixiv.net/ajax/illust/%d/pages?lang=zh", illustId))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resultError, ok := result.Error.(bool); ok && resultError {
		return nil, errors.Wrap(OriginUrlNoFound, result.Message)
	}

	var urls []string

	for _, v := range result.Body {
		urls = append(urls, v.Urls.Original)
	}

	if len(urls) <= 0 {
		return nil, OriginUrlNoFound
	}

	return urls, nil
}

func (client *Pixiv) RankImageUrls(options ...GetRankOptions) ([]Image, error) {
	p := 0
	imagesQuantity := 0
	rank := 0

	params := GetRankParams{
		Date: time.Time{},
	}

	for _, p := range options {
		p(&params)
	}

	if params.MaxQuantity != 0 {
		params.MaxPage = 0
	}

	var Images []Image

	var baseUrl string

	//https://www.pixiv.net/ranking.php?mode=daily&content=illust&date=20220825&p=2&format=json

	if params.Date.IsZero() {
		baseUrl = "https://www.pixiv.net/ranking.php?p=%d&format=json&content=illust&mode=daily&lang=zh"
	} else {
		baseUrl = "https://www.pixiv.net/ranking.php?p=%d&format=json&content=illust&lang=zh&mode=daily&date=" + params.Date.Format("20060102")
	}

	GetList := func(p int) (*RankListResult, error) {
		response, err := client.httpClient.R().Get(fmt.Sprintf(baseUrl, p))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if response.StatusCode != 200 {
			return nil, RankNotFound
		}

		result := RankListResult{}
		_ = json.Unmarshal(response.Bytes(), &result)
		if result.Error != "" {
			return nil, RankNotFound
		}
		return &result, nil
	}

	GetAllImagesUrl := func(content RankListContent, rank int) (*Image, error) {
		var imageUrls []ImageUrl
		urls, err := client.GetOriginUrl(content.IllustId)
		if err != nil {
			return nil, err
		}
		for _, url := range urls {
			urlArr := strings.Split(url, "/")
			imageUrls = append(imageUrls, ImageUrl{
				Url:  url,
				Name: urlArr[len(urlArr)-1],
			})
		}

		return &Image{
			Rank:     rank,
			Title:    content.Title,
			IllustId: content.IllustId,
			Urls:     imageUrls,
		}, nil
	}

	var rankContents []RankListContent
	for {
		p++
		result, err := GetList(p)
		if err != nil {
			return nil, err
		}

		if nextBool, ok := result.Next.(bool); ok && !nextBool {
			break
		}

		imagesQuantity += len(result.Contents)

		rankContents = slice.Concat(rankContents, result.Contents)

		if params.MaxPage != 0 && p >= params.MaxPage {
			break
		}

		if params.MaxQuantity != 0 && imagesQuantity >= params.MaxQuantity {
			break
		}
	}

	gPool := gpool.NewGPool(50)

	for _, content := range rankContents {
		rank++
		gPool.Add(1)
		go func(content RankListContent, rank int) {
			defer gPool.Done()
			_img, err := GetAllImagesUrl(content, rank)
			if err != nil {
				logger.Logger.Error("获取pixiv图片url出现错误:" + err.Error())
			}
			Images = append(Images, *_img)
		}(content, rank)
	}
	gPool.Wait()

	sort.Slice(Images, func(i, j int) bool {
		return Images[i].Rank < Images[j].Rank
	})

	if params.MaxQuantity != 0 && params.MaxQuantity < len(Images) {
		return Images[:params.MaxQuantity], nil
	}

	return Images, nil
}

func (client *Pixiv) GetImageContent(url, illustId string) ([]byte, error) {
	response, err := client.httpClient.R().
		SetHeader("referer", "https://www.pixiv.net/member_illust.php?mode=medium&illust_id="+illustId).
		Get(url)
	if err != nil {
		return nil, errors.WithMessage(ImageDownloadError, err.Error())
	}
	return response.Bytes(), nil
}
