package onedrive

import (
	"fmt"
	"os"
	"pixivImages/config"
	"pixivImages/database"
	"testing"
)

var drive *OneDrive

func init() {
	config.LoadConfig()
	database.InitRedis()
	drive = NewOneDrive()
}

func TestOneDrive_UploadFile(t *testing.T) {
	filePath := "C:\\Users\\niexiawei\\Downloads\\LegionPowerPlan.zip"
	f, _ := os.Open(filePath)
	result, err := drive.UploadFile("LegionPowerPlan.zip", "root", f)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", *result)
}

func TestOneDrive_UploadBigFile(t *testing.T) {
	filePath := "C:\\Users\\niexiawei\\Downloads\\kcyb-master.zip"
	f, _ := os.Open(filePath)
	fileInfo, _ := f.Stat()
	result, err := drive.UploadBigFile("kcyb-master.zip", "root", f, fileInfo.Size())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", *result)
}

func TestOneDrive_GetDownloadUrl(t *testing.T) {
	url, err := drive.GetDownloadUrl("01OXMAQ6MKPGH2IB5R5FFKVMM2SIGI5FLD")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(url)
}
