## Introduction

mymon(MySQL-Monitor) -- MySQL数据库运行状态数据采集脚本，采集包括global status, global variables, slave status等。

由github.com/open-falcon/mymon 修改而来，
支持配置多个MySQL实例
## Installation

```bash
# set $GOPATH and $GOROOT

mkdir -p $GOPATH/src/github.com/coraldane
cd $GOPATH/src/github.com/coraldane
git clone https://github.com/coraldane/mymon.git

cd mymon
go get ./...
control build


```

## Configuration

```
{
	"log_level": "debug",
	"interval": 60,
	"connect_timeout": 5,
	"max_idle": 10,
	"falcon_client": "http://127.0.0.1:1988/v1/push",
	"defult_monitor_account":{
		"user": "viewer",
		"passwd": "mixuan"
	},
	"db_server_list": [
	{
		"endpoint": "mixuan1",
		"host": "idao.me",
		"port": 3306,
		"user": "viewer",
		"passwd": "mixuan"
	}
	]
}
```
defult_monitor_account:推荐用这种方式,所有mysql实例都配置相同的monitor账号和密码（做好权限控制）,之后mymon会自动发现本地的实例，并自动采集和上报数据.所有mymon的配置相同维护方便。
db_server_list:用于监控远程的数据实例
## MySQL metrics

请参考./metrics.txt，其中的内容，仅供参考，根据MySQL的版本、配置不同，采集到的metrics也有差别。


## Contributors
 - zhangxuliang mail: zhangxuliang.io@foxmail.com
 - libin  微信：libin_cc  邮件：libin_dba@xiaomi.com
 - coraldane 邮件：coraldane@163.com

