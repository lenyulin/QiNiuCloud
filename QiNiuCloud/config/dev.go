//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "remote_user:Pwd970203..@tcp(42.194.164.163:3306)/wedy",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
	MainOSS: MainOSSConfig{
		Cos: []CosConfig{
			CosConfig{
				BucketURL:  "https://eg7yrmglsgb1-1253517205.cos.ap-guangzhou.myqcloud.com",
				ServiceURL: "https://service.cos.myqcloud.com",
				SecretID:   "AKID4NUN60iDWbT8Zmna9Ucfgi8NZiiU4RaW",
				SecretKey:  "mKl2990xxmOn11YCmCPRK3Zub9UbYLvK",
			},
		},
	},
	BackupOSS: BackupOSSConfig{
		Tos: []TosConfig{
			TosConfig{
				Ak:       "AKLTMzk1MGU3NWJjODM1NDE1ZWFlOTM1MmE1MTkyODQ4YzE",
				Sk:       "WmpNNVltWXpNV0U1T1dWaE5ERXlNR0k0WldRNU9XVm1ORFF4Wm1JMU5UYw==",
				Endpoint: "https://tos-cn-shanghai.volces.com",
				Region:   "cn-shanghai",
			},
		},
	},
}

//
//var Config1 = config{
//	DB: DBConfig{
//		DSN: "xxxxx",
//	},
//	Redis: RedisConfig{
//		Addr: "xxxxx",
//	},
//	MainOSS: MainOSSConfig{
//		cos: []CosConfig{
//			CosConfig{
//				BucketURL:  "xxxxx",
//				ServiceURL: "xxxxx",
//				SecretID:   "xxxxx",
//				SecretKey:  "xxxxx",
//			},
//		},
//	},
//	BackupOSS: BackupOSSConfig{
//		tos: []TosConfig{
//			TosConfig{
//				Ak:       "xxxxx",
//				Sk:       "xxxxx",
//				Endpoint: "xxxxx",
//				Region:   "xxxxx",
//			},
//		},
//	},
//}
