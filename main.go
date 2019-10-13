package main

import (
	"encoding/json"
	"fmt"
	"goproj/api"
	"goproj/trans"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

func secFmt(x int) string {
	x /= 1000
	m := x / 60
	s := x % 60
	h := m / 60
	m = m % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

type mtask struct {
	Start, End       int
	FileName, Result string
	Code             int
}
type taskser []*mtask

func (t taskser) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t taskser) Less(i, j int) bool {
	return t[i].Start <= t[j].Start
}

// Swap swaps the elements with indexes i and j.
func (t taskser) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func main() {
	task, err := trans.NewTransTask(1, "/home/enh/share/Music/集团信息化会议录音2(1).MP3")
	grouctl := make(chan int, 2)
	defer task.Free()
	if err != nil {
		panic(err)
	}
	task.Audio2PCM16k()
	s, err := task.Split()
	if err != nil {
		panic(err)
	}
	//results:=make(chan task,len(s))
	wg := &sync.WaitGroup{}
	var t *mtask
	reg, _ := regexp.Compile(`(\d+)\-(\d+)`)
	fun := api.Audio2StrHelper(1537)
	tasks := make([]*mtask, 0)
	for _, file := range s {
		t = &mtask{FileName: file}
		tasks = append(tasks, t)
		go func(m *mtask, wg *sync.WaitGroup, gc chan int) {
			defer func() { <-gc }()
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					m.Code = -1
					switch err.(type) {
					case json.SyntaxError:
						m.Result = err.(*json.SyntaxError).Error()
					default:
						m.Result = "err.Error()"
					}
				}
			}()
			path := filepath.Join(task.Out, m.FileName)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			} else {
				regsult := reg.FindStringSubmatch(m.FileName)
				if 3 != len(regsult) {
					panic(fmt.Errorf("Name Eror %p", regsult))
				}
				m.Start, _ = strconv.Atoi(regsult[1])
				m.End, _ = strconv.Atoi(regsult[2])
				gc <- 1
				fmt.Println("runing 1.")
				x, err := fun(data)
				if err != nil {
					m.Code = -1
					m.Result = err.Error()
				} else {
					if 0 != int(x["err_no"].(float64)) {
						m.Code = int(x["err_no"].(float64))
						m.Result = x["err_msg"].(string)
					} else {
						m.Code = 0
						m.Result = x["result"].([]interface{})[0].(string)
					}
				}
			}

		}(t, wg, grouctl)
		wg.Add(1)
	}
	wg.Wait()
	sort.Sort(taskser(tasks))
	for _, i := range tasks {
		fmt.Printf("%s-%s :%s\n", secFmt(i.Start), secFmt(i.End), i.Result)
	}
}
