package rbac

/*
Cấu hình rbac
*/
func Init(configs ...Config) {
	if len(configs) == 0 {
		config = NewConfig()
	} else {
		config = configs[0]
	}
}

/*
Tạo cấu hình mặc định cho RBAC
*/
func NewConfig() Config {
	return Config{
		ForbidOverAllow:           true,
		RootAllow:                 false,
		AuthService:               "",
		MakeUnassignedRoutePublic: false,
	}
}
