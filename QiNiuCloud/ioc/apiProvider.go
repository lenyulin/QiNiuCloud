package ioc

import (
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager"
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/Providers"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"github.com/redis/go-redis/v9"
	ai3d "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ai3d/v20250513"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	"os"
)

func InitHunYuan3D() *ai3d.Client {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ai3d.tencentcloudapi.com"
	client, _ := ai3d.NewClient(credential, regions.Ashburn, cpf)
	return client
}

//NewHunyuanTo3D(l logger.LoggerV1, client *ai3d.Client) ProviderSpecificGenerator

func InitHunyuanTo3DProvider(l logger.LoggerV1, client *ai3d.Client) Providers.ProviderSpecificGenerator {
	return Providers.NewHunyuanTo3D(l, client)
}

// ModelAPIsProviderManager
func InitModelAPIsProviderManager(l logger.LoggerV1, providers Providers.ProviderSpecificGenerator) AsyncModelGenerationTaskManager.ModelAPIsProviderManager {
	var p []Providers.ProviderSpecificGenerator
	p = append(p, providers)
	return AsyncModelGenerationTaskManager.NewModelAPIsProviderManager(l, p)
}

//AsyncModelGenerationTaskManager.NewSyncModelGenerationTaskManager

func InitNewSyncModelGenerationTaskManager(l logger.LoggerV1, redis redis.Cmdable, providerManager AsyncModelGenerationTaskManager.ModelAPIsProviderManager) AsyncModelGenerationTaskManager.SyncModelGenerationTaskManager {
	return AsyncModelGenerationTaskManager.NewSyncModelGenerationTaskManager(l, redis, providerManager)
}
