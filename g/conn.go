package g

import (
	"fmt"

	"github.com/keimoon/gore"
	"github.com/projecteru/eru-agent/logs"
)

var Rds *gore.Pool

func InitialConn() {
	Rds = &gore.Pool{
		InitialConn: Config.Redis.Min,
		MaximumConn: Config.Redis.Max,
	}

	redisHost := fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port)
	if err := Rds.Dial(redisHost); err != nil {
		logs.Assert(err, "Redis init failed")
	}

	logs.Info("Global connections initiated")
}

func CloseConn() {
	Rds.Close()
	logs.Info("Global connections closed")
}

func GetRedisConn() *gore.Conn {
	conn, err := Rds.Acquire()
	if err != nil || conn == nil {
		logs.Assert(err, "Get redis conn")
	}
	return conn
}

func ReleaseRedisConn(conn *gore.Conn) {
	Rds.Release(conn)
}
