package ossx

import (
	"fmt"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

func InitTOS() *tos.ClientV2 {
	var (
		//ak       = os.Getenv("TOS_ACCESS_KEY")
		//sk       = os.Getenv("TOS_SECRET_KEY")
		ak       = "AKLTMzk1MGU3NWJjODM1NDE1ZWFlOTM1MmE1MTkyODQ4YzE"
		sk       = "WmpNNVltWXpNV0U1T1dWaE5ERXlNR0k0WldRNU9XVm1ORFF4Wm1JMU5UYw=="
		endpoint = "https://tos-cn-shanghai.volces.com"
		region   = "cn-shanghai"
	)
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
