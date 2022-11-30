package pixiv

import (
	"fmt"
	"pixivImages/config"
	"strings"
	"testing"
	"time"
)

var pixiv *Pixiv

func init() {
	config.LoadConfig()
	pixiv = NewPixiv()
}

func TestPixiv_RankImageUrls(t *testing.T) {
	date := time.Date(2022, 11, 15, 0, 0, 0, 0, time.Local)
	image, err := pixiv.RankImageUrls(WithRankListDate(date), WithRankQueryMaxQuantity(1))
	if err != nil {
		t.Error(err)
		return
	}
	for _, _image := range image {
		fmt.Printf("%+v\n", _image)
		for _, url := range _image.Urls {
			t.Log(url.GetPath())
		}
	}
}

func Test_urlToPath(t *testing.T) {
	url := "https://i.pximg.net/img-original/img/2022/11/14/08/09/40/102781051_p0.png"
	url = strings.Replace(url, "https://", "", -1)
	pathArr := strings.Split(url, "/")
	pathArr = pathArr[:len(pathArr)-1]
	t.Log(pathArr)
}
