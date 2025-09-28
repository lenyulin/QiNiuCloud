package ossx

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

func InitCOS(SecretID string, SecretKey string, BucketURL string, ServiceURL string) *cos.Client {
	u, _ := url.Parse(BucketURL)
	su, _ := url.Parse(ServiceURL)
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  SecretID,
			SecretKey: SecretKey,
		},
	})
	return client
}
