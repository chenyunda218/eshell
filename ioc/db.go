package ioc

import (
	"eshell/ioc/config"
	"fmt"
	"log/slog"

	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/ulog"
	"gorm.io/gorm"
)

var loggerSingleton = singleton.SingletonFactory(func() *slog.Logger {
	return slog.New(ulog.NewHandler(0))
	// return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}, singleton.Eager)

func GetLogHandler() *slog.Logger {
	return loggerSingleton.Get()
}

var jwtSignerSingleton = singleton.SingletonFactory(createJwtSignerSingleton, singleton.Eager)

func GetJwtSigner() *auth.JwtSigner {
	return jwtSignerSingleton.Get()
}

func createJwtSignerSingleton() *auth.JwtSigner {
	config := config.GetConfig()
	return auth.NewJwtSigner([]byte(config.GetString("jwt.secret")))
}

var dbSingleton = singleton.SingletonFactory(createDBInstance, singleton.Lazy)

func GetDBInstance() *gorm.DB {
	return dbSingleton.Get()
}

func createDBInstance() *gorm.DB {
	config := config.GetConfig()
	fmt.Println(config.GetString("database.username"))
	db, err := dao.NewConnection(
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.host"),
		config.GetString("database.port"),
		config.GetString("database.database"),
	)
	if err != nil {
		panic(err)
	}
	return db
}
