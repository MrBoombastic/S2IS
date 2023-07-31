package types

type UploadRequest struct {
	File      string `json:"file"`
	Extension string `json:"extension"`
}
type UploadResponse struct {
	Error    string `json:"error"`
	Filename string `json:"filename"`
}

type DeleteRequest struct {
	Filename string `json:"filename"`
}

type DeleteResponse struct {
	Error string `json:"error"`
}
