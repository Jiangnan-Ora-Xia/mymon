package g

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/toolkits/file"
	"github.com/toolkits/slice"
	"github.com/toolkits/sys"
	"io"
	"strconv"
	"strings"
)

func MysqlPorts() ([]int, error) {
	ports := []int{}
	bs, err := sys.CmdOutBytes("/bin/sh", "-c", "ss -t -l -n -p ")
	if err != nil {
		return ports, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))
	// ignore the first line
	line, err := file.ReadLine(reader)
	if err != nil {
		return ports, err
	}

	for {
		line, err = file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return ports, err
		}

		fields := strings.Fields(string(line))
		fieldsLen := len(fields)

		if fieldsLen < 4 {
			return ports, fmt.Errorf("output of format not supported")
		}

		portColumnIndex := 2
		if fieldsLen > 4 {
			portColumnIndex = 3
		}

		location := strings.LastIndex(fields[portColumnIndex], ":")
		port := fields[portColumnIndex][location+1:]
		users := fields[fieldsLen-1]
		if p, e := strconv.Atoi(port); e != nil {
			return ports, fmt.Errorf("parse port to int64 fail: %s", e.Error())
		} else if strings.Contains(users, "\"mysqld\"") {
			ports = append(ports, p)
		}

	}
	return slice.UniqueInt(ports), nil
}
