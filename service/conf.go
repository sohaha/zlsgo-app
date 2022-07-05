package service

import (
	"reflect"

	"zlsapp/conf"

	"github.com/sohaha/zlsgo/zerror"
	gconf "github.com/zlsgo/conf"
)

// Conf 配置项
type Conf struct {
	Base   conf.Base
	DB     conf.DB
	Wechat conf.Wechat
}

func InitConf(defaultSet []interface{}) *Conf {
	c := &Conf{}
	cfg := gconf.New(conf.FileName)

	for _, c := range defaultSet {
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

	zerror.Panic(cfg.Read())
	zerror.Panic(cfg.Unmarshal(&c))

	return c
}
