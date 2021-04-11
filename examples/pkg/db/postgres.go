package db

import (
	"easyquery/examples/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Postgres *gorm.DB

func InitDB() {
	dsn := DSN(config.GlobalProperties)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Logger = logger.Default.LogMode(logger.Info)
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour * 24)
	Postgres = db
}

func CloseDB() {
	db, err := Postgres.DB()
	if err != nil {
		panic(err)
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
}

func DSN(info *config.Properties) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		info.Database.Host, info.Database.Username, info.Database.Password, info.Database.Database, info.Database.Port)
}
