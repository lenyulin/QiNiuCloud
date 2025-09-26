package MultiOssManager

import "errors"

var (
	ErrUploadFailed    = errors.New("upload failed")
	ErrOssClientFailed = errors.New("oss client failed")
)
