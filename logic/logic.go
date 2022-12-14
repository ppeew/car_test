package logic

import (
	"bufio"
	"context"
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

func createRoutinue(ctx context.Context, wg *sync.WaitGroup) chan struct{} {
	timeout := make(chan struct{})
	go func() {
		defer close(timeout)
		defer wg.Done()
		select {
		case <-time.After(time.Second * 7):
			//获得超时信号
			timeout <- struct{}{}
		case <-ctx.Done():
			//获得退出信号
			return
		}
	}()
	return timeout
}

func RunTest() {
	count := 10 //题目数量
	var wg sync.WaitGroup
	score := 0
	tran := [4]string{"近光灯", "远光灯", "远近交替", "示宽灯"}
	used_que := make(map[int64]bool, 0)

	for i := 0; i < count; i++ {
		//生成题目
		var ls []model.Light
		total := d.Raw("SELECT * FROM lights").Scan(&ls).RowsAffected
		fmt.Printf("total: %v\n", total)
		r := rand.Int63n(total)
		for used_que[r] {
			//使用过了
			r = (r*143 + 21) % total
		}
		used_que[r] = true
		l := ls[r]
		fmt.Printf("\033[1;31;40mQuestion%v: %v\033[0m\n\033[1;35;48m(0:近光灯--1:远光灯--2:远近交替--3:示宽灯):\033[0m", i+1, l.Question)

		//每次生成题目调用一个协程，检查超时时间
		ctx, cancel := context.WithCancel(context.Background())
		wg.Add(1)
		timeout := createRoutinue(ctx, &wg)
		var ans int
		fmt.Scanln(&ans)
		select {
		case <-timeout:
			fmt.Printf("超时,正确答案是:%v\n", l.Answer)
		default:
			if tran[ans] == l.Answer {
				score += 100 / count
				fmt.Println("回答正确！")
			} else {
				fmt.Printf("回答错误，正确答案是:%v\n", l.Answer)
			}
		}
		//完成事件
		cancel()
		wg.Wait()
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
