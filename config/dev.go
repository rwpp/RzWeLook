//go:build !k8s

package config

var Config = WelookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/welook",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
