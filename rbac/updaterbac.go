package rbac

import (
	"github.com/casbin/casbin/v2"
	redisadapter "github.com/casbin/redis-adapter/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

/*
	Hàm ConnectDB nhận vào đường dẫn đến file config casbin (mặc định là "config/rbac_model.conf")

và database redis sẽ lưu các policy (mặc định là 3), trả về enforcer để sau này thêm, xóa, cập
nhật policy
*/
func ConnectCasbinDB(fileconfig string, database int) *casbin.Enforcer {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(viper.GetString("redis.network"), viper.GetString("redis.addr"), redis.DialDatabase(database))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	a, _ := redisadapter.NewAdapterWithPool(pool)

	e, _ := casbin.NewEnforcer(fileconfig, a)

	e.LoadPolicy()

	return e
}

/*
	Chuyển các role được hard code từ hệ thống cũ dựa trên path + method vào hệ thống mới

Lưu ý kết nối casbin database trước
*/
func AddOldRole(app *iris.Application, enforce *casbin.Enforcer) {
	routes := app.GetRoutes()
	for _, route := range routes {
		for i := range pathsRoles[route.Path][route.Method] {
			var role string
			switch i {
			case ADMIN:
				role = "ADMIN"
			case STUDENT:
				role = "STUDENT"
			case TRAINER:
				role = "TRAINER"
			case SALE:
				role = "SALE"
			case EMPLOYER:
				role = "EMPLOYER"
			case AUTHOR:
				role = "AUTHOR"
			case EDITOR:
				role = "EDITOR"
			case MAINTAINER:
				role = "MAINTAINER"
			}
			enforce.AddPolicy(role, viper.GetString("host")+route.Path, route.Method)
		}
	}

	enforce.SavePolicy()
}
