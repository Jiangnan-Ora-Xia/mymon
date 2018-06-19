package db

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/orm"
	"github.com/coraldane/mymon/g"
	"github.com/toolkits/logger"

	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	ormLock = new(sync.RWMutex)
)

func InitDatabase() {
	// set default database
	if g.Config().LogLevel == "debug" {
		orm.Debug = true
	}

	maxIdle := g.Config().MaxIdle
	t := time.NewTicker(time.Duration(time.Second * 60)).C
	for {
		defaultDb := false
		for _, server := range g.Config().DBServerList {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?loc=Local&parseTime=true",
				server.User, server.Passwd, server.Host, server.Port)
			fmt.Println(dsn)
			if defaultDb || orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxIdle) == nil {
				defaultDb = true
			}

			alias := g.Hostname(server) + fmt.Sprint(server.Port)
			orm.RegisterDataBase(alias, "mysql", dsn, maxIdle, maxIdle)
		}
		<-t
	}

}

func NewOrmWithAlias(alias string) (orm.Ormer, error) {
	ormLock.RLock()
	defer ormLock.RUnlock()
	_, err := orm.GetDB()
	if err != nil {
		return nil, err
	}
	_, err = orm.GetDB(alias)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	o.Using(alias)
	return o, nil
}

func QueryFirst(alias, strSql string, args ...interface{}) (orm.Params, error) {
	var maps []orm.Params
	ormer, err := NewOrmWithAlias(alias)
	if ormer == nil {
		return nil, err
	}
	num, err := ormer.Raw(strSql, args...).Values(&maps)
	if nil != err {
		logger.Errorln(num, err)
		return nil, err
	}
	if num > 0 {
		return maps[0], err
	}
	return nil, err
}

func QueryRows(alias, strSql string, args ...interface{}) ([]orm.Params, error) {
	var maps []orm.Params
	ormer, err := NewOrmWithAlias(alias)
	if ormer == nil {
		return nil, err
	}
	num, err := ormer.Raw(strSql, args...).Values(&maps)

	if nil != err {
		logger.Errorln(num, err)
		return nil, err
	}
	return maps, err
}
