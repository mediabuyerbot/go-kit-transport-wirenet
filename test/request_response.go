package test

type UpdateBalanceRequest struct {
	A int
	B int
}

type UpdateBalanceResponse struct {
	Sum int
	Err error
}

type UploadFileRequest struct {
	Filepath string
	Size     int64
	Name     string
}

type UploadFileResponse struct {
	Err error
}
