package pixiv

import (
	"strings"
	"time"
)

type ImageUrl struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (i ImageUrl) GetPath() []string {
	url := i.Url
	url = strings.Replace(url, "https://", "", -1)
	pathArr := strings.Split(url, "/")
	pathArr = pathArr[:len(pathArr)-1]
	return pathArr
}

type Image struct {
	Urls     []ImageUrl `json:"urls"`
	Rank     int        `json:"rank"`
	Name     string     `json:"name"`
	Title    string     `json:"title"`
	IllustId int        `json:"illustId"`
}

type GetRankParams struct {
	Date        time.Time
	MaxPage     int
	RankNums    int
	MaxQuantity int
}

type RankListContent struct {
	Title    string `json:"title"`
	IllustId int    `json:"illust_id"`
	Rank     int    `json:"-"`
}

type RankListResult struct {
	Contents []RankListContent `json:"contents"`
	Next     interface{}       `json:"next"`
	Error    string            `json:"error"`
}

type OriginUrlResult struct {
	Error   interface{} `json:"error"`
	Message string      `json:"message"`
	Body    []struct {
		Urls struct {
			Original string `json:"original"`
		} `json:"urls"`
	} `json:"body"`
}

type GetRankOptions func(p *GetRankParams)

func WithRankListDate(date time.Time) GetRankOptions {
	return func(p *GetRankParams) {
		p.Date = date
	}
}

func WithRankQueryMaxPage(page int) GetRankOptions {
	return func(p *GetRankParams) {
		p.MaxPage = page
	}
}

func WithRankQueryRankNums(nums int) GetRankOptions {
	return func(p *GetRankParams) {
		p.RankNums = nums
	}
}

func WithRankQueryMaxQuantity(quantity int) GetRankOptions {
	return func(p *GetRankParams) {
		p.MaxQuantity = quantity
	}
}
