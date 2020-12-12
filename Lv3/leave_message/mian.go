package main

import (
	"leave_message/cmd"
	"leave_message/dao"
)

func main() {
	dao.MysqlInit()
	cmd.Entrance()
}
