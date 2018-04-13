package initializers

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var (
	DB *gorm.DB
	dbErr error
)

func logger() *log.Logger {
	fileName := beego.AppPath + "/logs/"+ beego.BConfig.RunMode + "_sql.log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("open file error")
	}

	return log.New(logFile, "\r\n", 0)
}

func init() {
	beego.Info("start initializer global variable gorm db")

	driverName := beego.AppConfig.String("default_database_driverName")
	dataSource := beego.AppConfig.String("default_database_dataSource")
	DB, dbErr = gorm.Open(driverName, dataSource)
	beego.Info("initializer global variable gorm db err =", dbErr)

	if dbErr == nil {
		maxIdle, _ := beego.AppConfig.Int("default_database_maxIdle")
		maxConn, _ := beego.AppConfig.Int("default_database_maxConn")

		DB.DB().SetMaxIdleConns(maxIdle)
		DB.DB().SetMaxOpenConns(maxConn)

		DB.LogMode(true)
		DB.SetLogger(gorm.Logger{logger()})
	}
	// defer DB.Close()
}
