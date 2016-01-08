package defines

type RedisConfig struct {
	Host string
	Port int
	Min  int
	Max  int
}

type SSHConfig struct {
	Bind    string
	HostKey string
	PrivKey string

	Redis RedisConfig
}
