package onedrive

type PathType string
type UploadType string

const (
	PathTypeFile         PathType   = "file"
	PathTypeFolder       PathType   = "folder"
	UploadTypeUploadFile UploadType = "upload"
	UploadTypeCreateFile UploadType = "create"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type CreateUploadSessionResponse struct {
	Url string `json:"uploadUrl"`
}

type UploadFileResult struct {
	Name       string     `json:"name"`
	Id         string     `json:"id"`
	UploadType UploadType `json:"-"`
}

type GetDownloadUrlResult struct {
	Id  string `json:"id"`
	Url string `json:"@microsoft.graph.downloadUrl"`
}
