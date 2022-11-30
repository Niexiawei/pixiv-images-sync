package onedrive

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"io"
	"pixivImages/app/service/ms_graph"
	"time"
)

const (
	BaseUrl          = "https://graph.microsoft.com/v1.0"
	BigFileBytesSize = 32 * 327680
)

var (
	UploadFileError  = errors.New("文件上传失败")
	GetDownloadError = errors.New("获取下载链接失败")
)

type OneDrive struct {
	HttpClient                   *req.Client
	authorizationRefreshStopChan chan struct{}
}

func NewOneDrive() *OneDrive {
	drive := &OneDrive{
		HttpClient: req.C().SetBaseURL(BaseUrl).
			SetCommonHeaders(map[string]string{
				"Authorization": "Bearer " + ms_graph.GetGraphToken(),
				"Content-type":  "application/json",
			}),
		authorizationRefreshStopChan: make(chan struct{}, 1),
	}

	go func() {
		tick := time.NewTicker(5 * time.Minute)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				drive.UploadAuthorization()
			case <-drive.authorizationRefreshStopChan:
				return
			}
		}
	}()

	return drive
}

func (drive *OneDrive) StopRefreshAuthorization() {
	drive.authorizationRefreshStopChan <- struct{}{}
}

func (drive *OneDrive) UploadAuthorization() {
	drive.HttpClient.SetCommonHeader("Authorization", "Bearer "+ms_graph.GetGraphToken())
}

func (drive *OneDrive) GetDownloadUrl(itemId string) (string, error) {
	errResp, result := &ErrorResponse{}, &GetDownloadUrlResult{}
	url := fmt.Sprintf("/me/drive/items/%s?select=id,@microsoft.graph.downloadUrl", itemId)
	response, err := drive.HttpClient.R().SetError(errResp).SetResult(result).Get(url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !response.IsSuccess() {
		return "", errors.WithMessage(GetDownloadError, errResp.Error.Message)
	}

	return result.Url, nil
}

func (drive *OneDrive) UploadBigFile(fileName, path string, reader io.Reader, size int64) (*UploadFileResult, error) {
	uploadResult, err := drive.UploadFile(fileName, path, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uploadUrl, err := drive.CreateUploadSession(uploadResult.Id, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer drive.CloseOverSizeUrl(uploadUrl)
	if _, err = drive.UploadOverSizedFiles(uploadUrl, size, reader); err != nil {
		drive.DeleteFile(uploadResult.Id)
		return nil, errors.WithStack(err)
	}
	return uploadResult, nil
}

func (drive *OneDrive) CloseOverSizeUrl(url string) bool {
	response, err := drive.HttpClient.R().Delete(url)
	if err != nil {
		return false
	}

	if !response.IsSuccess() {
		return false
	}

	return true
}

func (drive *OneDrive) UploadFile(fileName, path string, reader io.Reader) (*UploadFileResult, error) {
	errResp, result := &ErrorResponse{}, &UploadFileResult{}
	var url string
	if path == "root" {
		url = fmt.Sprintf("/me/drive/items/root:/%s:/content", fileName)
	} else {
		url = fmt.Sprintf("/me/drive/items/root:/%s/%s:/content", path, fileName)
	}
	response, err := drive.HttpClient.R().
		SetHeader("Content-type", "text/plain").
		SetBody(reader).
		SetError(errResp).
		SetResult(result).
		Put(url)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return nil, errors.WithMessage(UploadFileError, errResp.Error.Message)
	}

	if response.StatusCode == 200 {
		result.UploadType = UploadTypeUploadFile
	} else {
		result.UploadType = UploadTypeCreateFile
	}

	return result, nil
}

func (drive *OneDrive) DeleteFile(id string) bool {
	response, err := drive.HttpClient.R().Delete("/me/drive/items/" + id)
	if err != nil {
		return false
	}
	if !response.IsSuccess() {
		return false
	}
	return true
}

func (drive *OneDrive) CreateUploadSession(path string, deferCommit bool) (string, error) {
	errResp, result := &ErrorResponse{}, &CreateUploadSessionResponse{}
	url := fmt.Sprintf("/me/drive/items/%s/microsoft.graph.createUploadSession", path)
	context2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := drive.HttpClient.R().
		SetContext(context2).
		SetRetryCount(3).
		SetBody(map[string]interface{}{
			"deferCommit": deferCommit,
		}).
		SetError(errResp).
		SetResult(result).
		Post(url)

	if err != nil {
		return "", errors.WithStack(err)
	}

	if !response.IsSuccess() {
		return "", errors.New(errResp.Error.Message)
	}

	return result.Url, nil
}

func (drive *OneDrive) UploadOverSizedFiles(url string, size int64, reader io.Reader) (writerLen int, err error) {
	errResp := &ErrorResponse{}
	var currenUploadByte int
	upload := func(fileByte []byte) (int, error) {
		currByteLen := len(fileByte)

		var (
			response *req.Response
			errHttp  error
		)
		response, errHttp = drive.HttpClient.R().
			SetHeaders(map[string]string{
				"Content-Length": fmt.Sprintf("%d", currByteLen),
				"Content-Range":  fmt.Sprintf("bytes %d-%d/%d", currenUploadByte, (currenUploadByte+currByteLen)-1, size),
			}).
			SetBodyBytes(fileByte).
			SetError(errResp).
			Put(url)

		if errHttp != nil {
			return 0, errors.WithStack(errHttp)
		}

		if !response.IsSuccess() {
			return 0, errors.Wrapf(UploadFileError, "%s:%s", errResp.Error.Code, errResp.Error.Message)
		}

		return currByteLen, nil
	}
	ff := make([]byte, BigFileBytesSize)

	for {
		nn, errR := reader.Read(ff)
		if nn > 0 {
			uploadByte, errW := upload(ff[0:nn])
			if errW != nil {
				return 0, errW
			}
			currenUploadByte += uploadByte
		}
		if errR != nil {
			if errR != io.EOF {
				err = errR
			}
			break
		}
	}
	return currenUploadByte, err
}
