package common

import (
	"fmt"
	"testing"
)

func TestGetResultRow(t *testing.T) {
	conn, err := NewMysqlConn()
	if err != nil {
		panic(err)
	}

    sql := "select * from product where id=?"
	row, err := conn.Query(sql, 1)
	if err != nil {
		panic(err)
	}

	resultRow := GetResultRow(row)
	fmt.Println(resultRow)
}

func TestGetResultRows(t *testing.T) {
	conn, err := NewMysqlConn()
	if err != nil {
		panic(err)
	}

	sql := "select * from product where id=?"
	row, err := conn.Query(sql, 1)
	if err != nil {
		panic(err)
	}

	resultRows := GetResultRows(row)

	fmt.Println(resultRows)
}