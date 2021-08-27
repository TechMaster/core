package db

import (
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
)

var DB *pg.DB

func ConnectPostgresqlDB(){
	DB = pg.Connect(&pg.Options{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Database: viper.GetString("database.database"),
		Addr:     viper.GetString("database.address"),
	})
}
