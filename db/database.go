package db

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/orm"
	log "github.com/cihub/seelog"
	"github.com/coraldane/mymon/g"

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
	defaultDb := false
	maxIdle := g.Config().MaxIdle
	t := time.NewTicker(time.Duration(time.Second * 60)).C
	for {
		//auto discovery local db
		ports, err := g.MysqlPorts()
		if err != nil {
			log.Error(err)
		}
		for _, port := range ports {
			exist := false
			for _, server := range g.Config().DBServerList {
				if port == server.Port && server.Host == "127.0.0.1" {
					exist = true
				}
			}
			if !exist {
				discover_db := g.DBServer{
					Endpoint: g.Hostname(""),
					Host:     "127.0.0.1",
					Port:     port,
					User:     g.Config().Default_Monitor_Account.User,
					Passwd:   g.Config().Default_Monitor_Account.Passwd,
				}
				g.Config().DBServerList = append(g.Config().DBServerList, &discover_db)
				log.Info("auto discover local db port:", port)
			}
		}
		for _, server := range g.Config().DBServerList {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?loc=Local&parseTime=true",
				server.User, server.Passwd, server.Host, server.Port)
			if defaultDb || orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxIdle) == nil {
				defaultDb = true
			}

			alias := g.Hostname(server.Endpoint) + fmt.Sprint(server.Port)
			_, err = orm.GetDB(alias)
			if err != nil {
				orm.RegisterDataBase(alias, "mysql", dsn, maxIdle, maxIdle)
			}
			log.Info(server.String())
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
		log.Error(num, err)
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
		log.Error(num, err)
		return nil, err
	}
	return maps, err
}
