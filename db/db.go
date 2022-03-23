package db

import (
	"context"
	"fmt"

	"github.com/TechMaster/core/config"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/pgjson"
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

	pgjson.SetProvider(NewJSONProvider()) //Sử dụng goccy json

	if config.IsAppInDebugMode(){
		DB.AddQueryHook(dbLogger{}) //Log query to console
	}
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	bytes, _ := q.FormattedQuery()
	fmt.Println("After query :" + string(bytes))
	return nil
}
