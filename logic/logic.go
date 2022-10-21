package logic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ppeew/car_test/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var d *gorm.DB

func init() {
	dsn := "root:123456@tcp(localhost)/car?charset=utf8mb4&parseTime=True&loc=Local"
	d, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	d.AutoMigrate(&model.Light{})
}

var TIMEOUT chan string = make(chan string, 1) //是否超时
var mx sync.Mutex
var COUNT int = 10 //题目数量

func RunTest() {
	score := 0
	tran := [4]string{"近光灯", "远光灯", "远近交替", "示宽灯"}
	used_que := make(map[int64]bool, 0)

	t := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<-t.C
			mx.Lock()
			if len(TIMEOUT) == 0 {
				TIMEOUT <- "timeout"
				// fmt.Println("timeout++")
			}
			mx.Unlock()
		}
	}()

	for i := 0; i < COUNT; i++ {
		var ls []model.Light
		total := d.Raw("SELECT * FROM lights").Scan(&ls).RowsAffected
		r := rand.Int63n(total)
		for used_que[r] {
			//使用过了
			r = (r*143 + 21) % total
		}
		used_que[r] = true
		l := ls[r]
		fmt.Printf("\033[1;31;40mQuestion%v: %v\033[0m\n\033[1;35;48m(0:近光灯--1:远光灯--2:远近交替--3:示宽灯):\033[0m", i+1, l.Question)

		//设置定时器ticker配置
		t.Reset(5 * time.Second)
		//清空管道的剩余内容，定时器出现开始
		for len(TIMEOUT) == 1 {
			<-TIMEOUT
			// fmt.Println("timeput--")
		}

		var ans int
		fmt.Scanln(&ans)
		if tran[ans] == l.Answer && len(TIMEOUT) == 0 {
			score += 100 / COUNT
			fmt.Println("回答正确！")
		} else if len(TIMEOUT) == 1 {
			fmt.Printf("%v,正确答案是:%v\n", <-TIMEOUT, l.Answer) //取出管道
			// fmt.Println("timeout--")
		} else {
			fmt.Printf("回答错误，正确答案是:%v\n", l.Answer)
		}
	}
	fmt.Printf("\033[1;31;40m考试结束，你的成绩为:%v分\033[0m\n", score)
}

func InsertDB() {
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
		l := model.Light{Question: que, Answer: ans}
		d.Create(&l)
	}
	fmt.Printf("%v,数据全部插入完成\n", f.Name())
}
