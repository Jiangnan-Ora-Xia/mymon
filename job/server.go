package job

import (
	"fmt"
	"strconv"

	"github.com/coraldane/mymon/db"
	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/models"
)

func GlobalStatus(server *g.DBServer) ([]*models.MetaData, error) {
	return mysqlState(server, "SHOW /*!50001 GLOBAL */ STATUS")
}

func GlobalVariables(server *g.DBServer) ([]*models.MetaData, error) {
	return mysqlState(server, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

func mysqlState(server *g.DBServer, sql string) ([]*models.MetaData, error) {
        dbalias := g.Hostname(server)+fmt.Sprint(server.Port)
	rows, err := db.QueryRows(dbalias, sql)
	if err != nil {
		return nil, err
	}

	data := make([]*models.MetaData, len(rows))
	i := 0
	for _, row := range rows {
		key_ := fmt.Sprintf("%v", row["Variable_name"])
		v, _ := strconv.Atoi(fmt.Sprintf("%v", row["Value"]))

		data[i] = models.NewMetric(key_, server)
		data[i].SetValue(v)
		i++
	}
	return data[:i], nil
}
