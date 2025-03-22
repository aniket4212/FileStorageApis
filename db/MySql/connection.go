package MySql

import (
	"database/sql"
	"filestorage/config"
	"log"

	"github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func ConnectMysqlDB() {

	cfg := mysql.Config{
		User:   config.AppConfig.MysqlConf.Username,
		Passwd: config.AppConfig.MysqlConf.Password,
		Net:    config.AppConfig.MysqlConf.Net,
		Addr:   config.AppConfig.MysqlConf.Address,
		DBName: config.AppConfig.MysqlConf.DatabaseName,
	}

	// create a db connection pool
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Error while connecting to MySQL DB", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("Error while connecting to Mysql DB", err)
	}
	log.Println("Mysql Database Connected")
}

func CloseMysql() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Println("Error closing Mysql connection: ", err)
			return
		}
		log.Println("Closed Mysql connection Successfully!")
	}
}
