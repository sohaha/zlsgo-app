package example

import (
	"github.com/sohaha/zlsgo/zdi"
	"github.com/zlsgo/app_core/service"
)

// Conf 插件配置
type Conf struct {
	Dev  bool   `z:"dev"`
	Text string `z:"text"`
}

// ConfKey 配置文件key
func (Conf) ConfKey() string {
	return "demo"
}

// 可以禁止配置自动写入配置文件
// func (Conf) DisableWrite() bool {
// 	return true
// }

// 配置文件默认值
var conf = &Conf{
	Dev:  true,
	Text: "这是一个插件配置",
}

// 公开一个结构体给其他插件
type Share struct {
	Name string
}

// New 实例化插件
func New() (p *Plugin) {
	service.DefaultConf = append(service.DefaultConf, conf)

	return &Plugin{
		Pluginer: service.Pluginer{
			OnLoad: func(di zdi.Invoker) (any, error) {
				// 配置解析完成后执行

				var c *service.Conf
				_ = di.Resolve(&c)
				// 如果 conf 不是一个指针，那么这里需要使用 c.Unmarshal(conf.ConfKey(), conf)
				// if err := c.Unmarshal(conf.ConfKey(), &conf); err != nil {
				// 	return err
				// }
				p.log.Debug("插件配置：", conf)

				return func() *Share {
					// 注入依赖共享给其他插件使用
					return &Share{"demo"}
				}, nil
			},
			OnStart: func(di zdi.Invoker) error {
				// 全部插件加载完成后执行

				// 可以获取任意插件实例或者注入的依赖
				var share *Share
				if err := di.Resolve(&share); err != nil {
					return err
				}
				p.log.Debug("共享的数据：", share)

				return nil
			},
			OnDone: func(di zdi.Invoker) error {
				// 全部插件启动后执行
				return nil
			},
			OnStop: func(di zdi.Invoker) error {
				// 程序停止之前执行
				return nil
			},
			Service: &service.PluginService{
				Controllers: []service.Controller{&Index{}},
				Tasks: []service.Task{
					{
						Name: "demo task",
						Cron: "1 * * * * * *",
						Run: func() {
							p.log.Debug("定时执行任务")
						}},
				},
			},
		},
	}
}
