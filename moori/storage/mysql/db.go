package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moori/entity"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
type DB struct {
	config Config
	db     *gorm.DB
}

func (m *DB) Connect() *gorm.DB {
	return m.db
}
func New(config Config) *DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, oErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if oErr != nil {
		panic(oErr)
	}
	mySql, dErr := db.DB()
	if dErr != nil {
		fmt.Println("storage err: ", dErr)
	}

	mySql.SetMaxIdleConns(3)
	mySql.SetMaxOpenConns(100)
	mySql.SetConnMaxLifetime(time.Hour)

	return &DB{config: config, db: db}
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.Product{},
		&entity.Image{},
	)
	if err != nil {
		panic("AutoMigrate err: " + err.Error())
	}
}
