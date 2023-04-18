package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DBconn *gorm.DB

func Init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.%s_", os.Getenv("DB_SCHEMA"), os.Getenv("DB_TABLE_PREFIX")),
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Postgres Database Connected")
	conn.Logger = logger.Default.LogMode(logger.Info)

	sqlDB, _ := conn.DB()
	val, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONN"))
	if err != nil {
		val = 10
	}
	sqlDB.SetMaxIdleConns(val)

	val, err = strconv.Atoi(os.Getenv("MAX_OPEN_CONN"))
	if err != nil {
		val = 100
	}
	sqlDB.SetMaxOpenConns(val)

	val, err = strconv.Atoi(os.Getenv("MAX_CONN_LIFETIME"))
	if err != nil {
		val = 120
	}
	sqlDB.SetConnMaxLifetime(time.Duration(val) * time.Second)

	DBconn = conn
}
