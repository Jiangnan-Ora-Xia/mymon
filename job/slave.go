package job

import (
	"fmt"
	"strconv"

	log "github.com/toolkits/logger"

	"github.com/coraldane/mymon/db"
	"github.com/coraldane/mymon/g"
	"github.com/coraldane/mymon/models"
)

var SlaveStatusToSend = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_Log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func SlaveStatus(server *g.DBServer) ([]*models.MetaData, error) {
	isSlave := models.NewMetric("Is_slave", server)
	dbalias := g.Hostname(server) + fmt.Sprint(server.Port)
	log.Debug(dbalias)
	rows, err := db.QueryRows(dbalias, "SHOW SLAVE STATUS")
	if err != nil {
		return nil, err
	}

	// be master
	if 0 == len(rows) || rows[0] == nil {
		isSlave.SetValue(0)
		return []*models.MetaData{isSlave}, nil
	}

	// be slave
	isSlave.SetValue(1)

	data := make([]*models.MetaData, len(SlaveStatusToSend))
	for i, s := range SlaveStatusToSend {
		data[i] = models.NewMetric(s, server)
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].SetValue(1)
			for _, row := range rows {
				v := fmt.Sprintf("%v", row[s])
				if v != "Yes" {
					data[i].SetValue(0)
					break
				}
			}
		default:
			v, err := strconv.Atoi(fmt.Sprintf("%v", rows[0][s]))
			if err != nil {
				data[i].SetValue(-1)
			} else {
				data[i].SetValue(v)
			}
		}
	}
	return append(data, isSlave), nil
}
