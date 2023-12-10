package database

import (
	"database/sql"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	BlogDB    *gorm.DB
	BlogSqlDB *sql.DB
)

type mysqlConfig struct {
	Name      string `json:"name"`
	Host      string `json:"host"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Charset   string `json:"charset"`
	Collation string `json:"collation"`
	Port      string `json:"port"`
	Prefix    string `json:"prefix"`
}

func getConnection(dataBase string) (*gorm.DB, *sql.DB) {
	//dsn := "root:root@tcp(127.0.0.1:3306)/point_event?charset=utf8mb4&parseTime=True&loc=Local"

	goatGameDb := &mysqlConfig{}

	goatGameConf := viper.GetStringMapString("mysql." + dataBase)
	mapstructure.Decode(goatGameConf, &goatGameDb)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&collation=%s&parseTime=True&loc=Local",
		goatGameDb.Username,
		goatGameDb.Password,
		goatGameDb.Host,
		goatGameDb.Port,
		goatGameDb.Database,
		goatGameDb.Charset,
		goatGameDb.Collation,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}
	if viper.GetBool("debug") {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	sqlDb, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(2)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDb.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db, sqlDb
}

func InitMysql() {
	BlogDB, BlogSqlDB = getConnection("blog")
}
