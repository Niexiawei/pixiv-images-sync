package database

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"pixivImages/build"
	"pixivImages/config"
)

var _conn *gorm.DB

func GetMysqlDb() *gorm.DB {
	return _conn
}

func InitMysql() {
	var err error
	mysqlConfig := config.Get().Mysql
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Db,
	)

	_conn, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: func() logger.Interface {
			if build.Version == "prod" {
				return logger.Default.LogMode(logger.Error)
			} else {
				return logger.Default.LogMode(logger.Info)
			}
		}(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(errors.WithStack(err))
	}

	conn, err := _conn.DB()
	if err != nil {
		panic(errors.WithStack(err))
	}
	conn.SetMaxOpenConns(mysqlConfig.Pool)
	conn.SetMaxIdleConns(mysqlConfig.Pool)
}
