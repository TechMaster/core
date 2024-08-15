package rbac

import (
	"regexp"

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
func ConnectCasbinDB(fileconfig string, database int) (*casbin.Enforcer, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(viper.GetString("redis.network"), viper.GetString("redis.addr"), redis.DialDatabase(database))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	a, err := redisadapter.NewAdapterWithPool(pool)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(fileconfig, a)
	if err != nil {
		return nil, err
	}

	return e, e.LoadPolicy()
}

/*
	Chuyển các role được hard code từ hệ thống cũ dựa trên path + method vào hệ thống mới

Lưu ý kết nối casbin database trước
*/
func AddOldRole(app *iris.Application, enforce *casbin.Enforcer) (err error) {
	routes := app.GetRoutes()
	regxp := regexp.MustCompile(`:(([a-zA-Z_])*)`)

	for _, route := range routes {
		for i := range pathsRoles[regxp.ReplaceAllString(route.Path, "{$1}")].Roles {
			var role string
			for k, v := range Roles {
				if v == i {
					role = k
					break
				}
			}
			_, err = enforce.AddPolicy(role, viper.GetString("host")+route.Path, route.Method)
			if err != nil {
				return err
			}
		}
	}

	return enforce.SavePolicy()
}
