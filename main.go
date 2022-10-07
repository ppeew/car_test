package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type light struct {
	gorm.Model
	Question string `gorm:"type:TEXT;not null"`
	Answer   string `gotm:"not null"`
}

var COUNT int = 10

func main() {
	rand.Seed(time.Now().Unix())
	dsn := "root:123456@tcp(localhost)/car?charset=utf8mb4&parseTime=True&loc=Local"
	d, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	d.AutoMigrate(&light{})

	for {
		fmt.Printf("--1.开始云端考试\n--2.插入数据库题型\n(输入其他健可以退出)")
		fmt.Printf("请输入操作:")
		var need_quit bool = false
		var n string
		fmt.Scanln(&n)
		switch n {
		case "1":
			score := 0
			tran := [4]string{"近光灯", "远光灯", "远近交替", "示宽灯"}

			for i := 0; i < COUNT; i++ {
				var ls []light
				r := rand.Int63n(d.Raw("SELECT * FROM lights").Scan(&ls).RowsAffected)
				l := ls[r]
				fmt.Printf("Question%v: %v\n请输入您的答案(0:近光灯--1:远光灯--2:远近交替--3:示宽灯):", i+1, l.Question)
				var ans int
				fmt.Scanln(&ans)
				if tran[ans] == l.Answer {
					score += 100 / COUNT
					fmt.Println("回答正确！")
				} else {
					fmt.Printf("回答错误，正确答案是:%v\n", l.Answer)
				}
			}
			fmt.Printf("考试结束，你的成绩为:%v分\n", score)
		case "2":
			f, _ := os.OpenFile("light.txt", os.O_RDONLY, 0666)
			s := bufio.NewScanner(f)
			s.Split(bufio.ScanLines)
			for s.Scan() {
				s2 := bufio.NewScanner(strings.NewReader(s.Text()))
				s2.Split(bufio.ScanWords)
				s2.Scan()
				que := s2.Text()
				s2.Scan()
				ans := s2.Text()
				l := light{Question: que, Answer: ans}
				d.Create(&l)
			}
			fmt.Printf("%v,数据全部插入完成\n", f.Name())
		default:
			fmt.Println("bye!")
			need_quit = true
		}
		if need_quit {
			break
		}

		// var que, ans string
		// fmt.Println("请输入要插入的题目")
		// fmt.Scanln(&que)
		// fmt.Println("请输入对应的答案")
		// fmt.Scanln(&ans)
		// fmt.Printf("que: %v\n", que)
		// fmt.Printf("ans: %v\n", ans)
		// l := light{Question: que, Answer: ans}
		// d.Create(&l)
	}
}
