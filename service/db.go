package service

import (
	"errors"
	"strings"
	"zlsapp/internal/utils"

	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
	"github.com/zlsgo/zdb/driver"
	"github.com/zlsgo/zdb/driver/mysql"
	"github.com/zlsgo/zdb/driver/postgres"
	"github.com/zlsgo/zdb/driver/sqlite3"
)

func RegDB(c *Conf) *zdb.DB {
	var dbConf driver.IfeConfig
	d := strings.ToLower(c.DB.Driver)
	switch d {
	case "mysql":
		dbConf = &mysql.Config{
			Host:     c.DB.MySQL.Host,
			Port:     c.DB.MySQL.Port,
			User:     c.DB.MySQL.User,
			Password: c.DB.MySQL.Password,
			DBName:   c.DB.MySQL.DBName,
		}

	case "postgres":
		dbConf = &postgres.Config{
			Host:     c.DB.Postgres.Host,
			Port:     c.DB.Postgres.Port,
			User:     c.DB.Postgres.User,
			Password: c.DB.Postgres.Password,
			DBName:   c.DB.Postgres.DBName,
			SSLMode:  c.DB.Postgres.SSLMode,
		}
	case "sqlite":
		dbConf = &sqlite3.Config{
			File: c.DB.Sqlite.Path,
		}

	default:
		tip := "未知数据库驱动 " + d
		if d == "" {
			tip = "未配置数据库"
		}
		utils.Fatal(errors.New("初始化数据库失败: " + tip))
	}

	builder.DefaultDriver = dbConf.(driver.Dialect)

	db, err := zdb.New(dbConf)
	utils.Fatal(err)

	return db
}
