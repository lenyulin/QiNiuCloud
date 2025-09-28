package ossx

import (
	"fmt"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

func InitTOS(ak, sk, endpoint, region string) *tos.ClientV2 {
	credential := tos.NewStaticCredentials(ak, sk)
	// 可以通过 tos.WithRetry 的方式添加重试次数
	client, err := tos.NewClientV2(endpoint, tos.WithCredentials(credential), tos.WithRegion(region), tos.WithMaxRetryCount(3))
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer client.Close()
	return client
	// 使用结束后，关闭 client
	//client.Close()
}
