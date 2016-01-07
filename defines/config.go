package defines

type RedisConfig struct {
	Host string
	Port int
	Min  int
	Max  int
}

type SSHConfig struct {
	Bind string
	Key  string

	Redis RedisConfig
}
