package session

import (
	"time"

	"github.com/TechMaster/core/sessions"
	redis_session "github.com/TechMaster/core/sessions/sessiondb/redis"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	SESSION_COOKIE = "mycookiesession"
	SESS_USER      = "authenticate"
	AUTHINFO       = "authinfo"
)

//Các biến dùng chung trong packge
var Sess *sessions.Sessions         //Cấu hình Session Manager
var redisDB *redis_session.Database //Đây là một wrapper nối xuống Redis của Iris
var RedisClient *redis.Client       //Đây là redis client trực tiếp nối xuống Redis db không qua Iris
var expires = time.Hour * 720         //Thời gian mà 1 session sẽ hết hạn và bị xoá khỏi Redis

/* Khởi tạo In Memory Session, không kết nối vào Redis hay bất kỳ CSDL nào
Dùng trong ứng dụng đơn lẻ
*/
func init() {
	Sess = sessions.New(sessions.Config{
		Cookie:       SESSION_COOKIE,
		AllowReclaim: true,
		Expires:      expires,
		CookieSecureTLS: true,
	})
}

/*
Khi có nhiều web site dùng chung Session, cần lưu Session vào Redis database
Hàm này thay thế cho InitSession() vì có thể trong tương lai có thêm lựa chọn
lưu session vào MySQL, MongoDB hoặc Postgresql
*/
func InitRedisSession() *redis_session.Database {
	redisDB = redis_session.New(redis_session.Config{
		Network:   viper.GetString("redis.network"),
		Addr:      viper.GetString("redis.addr"),
		Password:  viper.GetString("redis.password"),
		Database:  viper.GetString("redis.database"),
		MaxActive: viper.GetInt("redis.max_active"),
		Timeout:   time.Duration(viper.GetInt("redis.idle_timeout")) * time.Minute,
		Prefix:    viper.GetString("redis.prefix"),
		Driver:    redis_session.GoRedis(),
	})

	Sess.UseDatabase(redisDB)

	RedisClient = redis.NewClient(&redis.Options{
		Network:  viper.GetString("redis.network"),
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       0,
	})

	return redisDB
}
