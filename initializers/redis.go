package initializers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
)

var (
	Redis redis.Conn
	redisErr error
)

func init() {
	beego.Info("start initializer global variable redis")

	network := beego.AppConfig.String("redis_network")
	address := beego.AppConfig.String("redis_address")
	db := beego.AppConfig.String("redis_db")

	Redis, redisErr = redis.Dial(network, address)
	beego.Info("initializer global variable redis err =", redisErr)

	if len(db) != 0 && redisErr == nil {
		_, err := Redis.Do("SELECT", db)
		beego.Info("initializer global variable redis success after, select db err =", err)
	}
	// defer Redis.Close()
}
