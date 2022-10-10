package main

import (
	"fmt"
	"math/rand"
	"time"

	"car_test/logic"
)

var need_quit bool = false                     //是否退出程序

func main() {

	rand.Seed(time.Now().Unix())
	for {
		fmt.Printf("--1.开始考试\n--2.插入数据库题型\n(输入其他健可以退出)")
		fmt.Printf("请输入操作:")
		var n string
		fmt.Scanln(&n)
		switch n {
		case "1":
			logic.RunTest()
		case "2":
			logic.InsertDB()
		default:
			fmt.Println("bye!")
			need_quit = true
		}
		if need_quit {
			break
		}

	}
}
