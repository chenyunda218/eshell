package ioc

import (
	"eshell/client"
	"eshell/ioc/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/universalmacro/common/singleton"
)

func GetEsClient() *client.Client {
	return esSign.Get()
}

var esSign = singleton.SingletonFactory(func() *client.Client {
	config := config.GetConfig()
	c, _ := client.New(elasticsearch.Config{
		CertificateFingerprint: config.GetString("es.certificateFingerprint"),
		Username:               config.GetString("es.username"),
		Password:               config.GetString("es.password"),
		Addresses:              config.GetStringSlice("es.addresses"),
	},
		config.GetString("es.username"),
		config.GetString("es.password"),
		config.GetString("es.kibanaUrl"),
	)
	return c
}, singleton.Eager)
