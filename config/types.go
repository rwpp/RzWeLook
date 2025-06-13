package config

type WelookConfig struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	DB       int
	Password string
}
