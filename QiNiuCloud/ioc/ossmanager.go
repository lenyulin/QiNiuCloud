package ioc

import (
	"QiNiuCloud/QiNiuCloud/config"
	"QiNiuCloud/QiNiuCloud/ioc/ossx"
	"QiNiuCloud/QiNiuCloud/pkg/MultiOssManager"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	oss "QiNiuCloud/QiNiuCloud/pkg/ossx"
	cos2 "QiNiuCloud/QiNiuCloud/pkg/ossx/cos"
	tos2 "QiNiuCloud/QiNiuCloud/pkg/ossx/tos"
)

func InitMultiOssManager(l logger.LoggerV1) MultiOssManager.MultiOssManager {
	return MultiOssManager.NewOSSManager(l, createMainOSS(), creatCosClient())
}
func createMainOSS() []oss.OSSHandler {
	var mainOss []oss.OSSHandler
	for _, cosConf := range config.Config.MainOSS.Cos {
		client := ossx.InitCOS(cosConf.SecretID, cosConf.SecretKey, cosConf.BucketURL, cosConf.ServiceURL)
		ossHdl := cos2.NewCOSHandler(client)
		mainOss = append(mainOss, ossHdl)
	}
	for _, tosConf := range config.Config.MainOSS.Tos {
		client := ossx.InitTOS(tosConf.Ak, tosConf.Sk, tosConf.Endpoint, tosConf.Region)
		ossHdl := tos2.NewTOSHandler(client)
		mainOss = append(mainOss, ossHdl)
	}
	return mainOss
}
func creatCosClient() []oss.OSSHandler {
	var backupOss []oss.OSSHandler
	for _, cosConf := range config.Config.BackupOSS.Cos {
		client := ossx.InitCOS(cosConf.SecretID, cosConf.SecretKey, cosConf.BucketURL, cosConf.ServiceURL)
		ossHdl := cos2.NewCOSHandler(client)
		backupOss = append(backupOss, ossHdl)
	}
	for _, tosConf := range config.Config.BackupOSS.Tos {
		client := ossx.InitTOS(tosConf.Ak, tosConf.Sk, tosConf.Endpoint, tosConf.Region)
		ossHdl := tos2.NewTOSHandler(client)
		backupOss = append(backupOss, ossHdl)
	}
	return backupOss
}
