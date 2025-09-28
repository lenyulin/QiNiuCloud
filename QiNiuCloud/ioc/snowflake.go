package ioc

import "QiNiuCloud/QiNiuCloud/pkg/snowflake"

func InitSnowflake() *snowflake.Snowflake {
	c, _ := snowflake.NewSnowflake(1, 1)
	return c
}
