package config

type config struct {
	DB        DBConfig
	Redis     RedisConfig
	MainOSS   MainOSSConfig
	BackupOSS BackupOSSConfig
}
type DBConfig struct {
	DSN string
}
type RedisConfig struct {
	Addr string
}
type MainOSSConfig struct {
	Cos []CosConfig
	Tos []TosConfig
}
type BackupOSSConfig struct {
	Cos []CosConfig
	Tos []TosConfig
}
type CosConfig struct {
	BucketURL  string
	ServiceURL string
	SecretID   string
	SecretKey  string
}
type TosConfig struct {
	Ak       string
	Sk       string
	Endpoint string
	Region   string
}
