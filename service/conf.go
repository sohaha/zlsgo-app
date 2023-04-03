package service

import (
	"reflect"

	"zlsapp/conf"
	"zlsapp/internal/utils"

	"github.com/spf13/viper"
	gconf "github.com/zlsgo/conf"
)

// Conf 配置项
type Conf struct {
	cfg *gconf.Confhub

	Base conf.Base

	DB conf.DB
}

func RegConf() *Conf {
	cfg := gconf.New(conf.FileName)
	c := &Conf{cfg: cfg}

	for _, c := range conf.DefaultConf {
		m, t := make(map[string]interface{}), reflect.TypeOf(c)

		isPrt := t.Kind() == reflect.Ptr
		if isPrt {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			continue
		}

		v := reflect.ValueOf(c)
		if isPrt {
			v = v.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			value, field := v.Field(i), t.Field(i)
			if value.IsZero() {
				continue
			}

			m[field.Name] = v.Field(i).Interface()
		}

		cfg.SetDefault(t.Name(), m)
	}

	utils.Fatal(cfg.Read())
	utils.Fatal(cfg.Unmarshal(&c))

	return c
}

func (c *Conf) Core() *viper.Viper {
	return c.cfg.Core
}
