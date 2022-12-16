package mysql

import (
	"awesomeProject/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConf) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connnect mysql failed, err: %v\n", err)
		return
	}
	fmt.Printf("connect mysql success....")
	return
}

func Close() {
	_ = db.Close()
}
