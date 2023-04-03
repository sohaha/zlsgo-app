package conf

type Base struct {
	Name        string `mapstructure:"name"`         // 项目名称
	Debug       bool   `mapstructure:"debug"`        // 开启全局调试模式
	LogDir      string `mapstructure:"log_dir"`      // 日志目录
	LogPosition bool   `mapstructure:"log_position"` // 调试下打印日志显示输出位置
	Port        string // 项目端口
	Pprof       bool   // 开启 pprof
	PprofToken  string // pprof Token
}

const (
	// FileName 配置文件名
	FileName = "conf"
	// LogPrefix 日志前缀
	LogPrefix = ""
	// AppName 项目名称
	AppName = "ZlsApp"
	// AppVersion 项目版本
	AppVersion = "1.0.0"

	debug = true
)

var (
	DefaultConf []interface{}
)

func init() {
	DefaultConf = append(DefaultConf, Base{
		Name:  AppName,
		Debug: debug,
		Port:  "8181",
	})
}
