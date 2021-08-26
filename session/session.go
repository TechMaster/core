package session

import (
	"time"

	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"github.com/spf13/viper"
)

const (
	SESSION_COOKIE = "sessid"
	SESS_AUTH      = "authenticated"
	SESS_USER      = "user"
	AUTHINFO       = "authinfo"
)

/*
Cấu hình Session Manager
*/
var Sess *sessions.Sessions

/* Khởi tạo In Memory Session, không kết nối vào Redis hay bất kỳ CSDL nào
Dùng trong ứng dụng đơn lẻ
*/
func init() {
	Sess = sessions.New(sessions.Config{
		Cookie:       SESSION_COOKIE,
		AllowReclaim: true,
		Expires:      time.Hour * 48, /*Có giá trị trong 2 ngày*/
	})
}

/*
Khi có nhiều web site dùng chung Session, cần lưu Session vào Redis database
Hàm này thay thế cho InitSession() vì có thể trong tương lai có thêm lựa chọn
lưu session vào MySQL, MongoDB hoặc Postgresql
*/
func InitRedisSession() *redis.Database {
	redisDB := redis.New(redis.Config{
		Network:   viper.GetString("redis.network"),
		Addr:      viper.GetString("redis.addr"),
		Password:  viper.GetString("redis.password"),
		Database:  viper.GetString("redis.database"),
		MaxActive: viper.GetInt("redis.max_active"),
		Timeout:   time.Duration(viper.GetInt("redis.idle_timeout")) * time.Minute,
		Prefix:    viper.GetString("redis.prefix"),
		Driver:    redis.GoRedis(),
	})

	Sess.UseDatabase(redisDB)
	return redisDB
}
