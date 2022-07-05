package conf

type (
	DB struct {
		Driver   string
		MySQL    mysql
		Postgres postgres
		Sqlite   sqlite
	}

	mysql struct {
		Host       string
		Port       int
		User       string
		Password   string
		DBName     string
		Parameters string
	}
	postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	sqlite struct {
		Path string
	}
)

func init() {
	DefaultSet = append(DefaultSet, DB{
		Driver: "mysql",
		MySQL: mysql{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "666666",
			DBName:   "zls",
		},
	})
}
