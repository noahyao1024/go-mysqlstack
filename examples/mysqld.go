/*
 * go-mysqlstack
 * xelabs.org
 *
 * Copyright (c) XeLabs
 * GPL License
 *
 */

package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/xelabs/go-mysqlstack/driver"
	querypb "github.com/xelabs/go-mysqlstack/sqlparser/depends/query"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
	"github.com/xelabs/go-mysqlstack/xlog"
)

func main() {
	port := 4407
	if len(os.Args) >= 2 {
		port, _ = strconv.Atoi(os.Args[1])
	}

	result1 := &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "id",
				Type: querypb.Type_INT32,
			},
			{
				Name: "name",
				Type: querypb.Type_VARCHAR,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte("10")),
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("nice name")),
			},
		},
	}

	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	th := driver.NewTestHandler(log)
	th.AddQuery("SELECT * FROM MOCK", result1)
	/*
		th.AddQuery("insert into `mysql_stack_mock_test` (`token`) values ('ruansishi')", &sqltypes.Result{
			InsertID: 123,
		})
	*/

	// For initialization
	/*
		mysql> select version();
		+------------+
		| version()  |
		+------------+
		| 5.7.24-log |
		+------------+
		1 row in set (0.00 sec)
	*/

	th.AddQuery("select version()", &sqltypes.Result{
		Fields: []*querypb.Field{
			{
				Name: "version()",
				Type: querypb.Type_VARCHAR,
			},
		},
		Rows: [][]sqltypes.Value{
			{
				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte("5.7.24-log")),
			},
		},
	})

	for _, sql := range []string{"start transaction", "commit"} {
		th.AddQuery(sql, &sqltypes.Result{
			Rows: [][]sqltypes.Value{},
		})
	}

	mysqld, err := driver.MockMysqlServerWithPort(log, port, th)
	if err != nil {
		log.Panic("mysqld.start.error:%+v", err)
	}
	defer mysqld.Close()
	log.Info("mysqld.server.start.address[%v]", mysqld.Addr())

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
