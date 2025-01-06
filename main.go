package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	rawfile  = "data.csv"
	goodfile = "good.csv"
	badfile  = "bad.csv"
	raw      *os.File
	good     *os.File
	bad      *os.File
	err      error
	pool     *sync.Pool
)

func Task() {
	if raw, err = os.OpenFile(rawfile, os.O_RDWR, 0644); err != nil {
		panic(err)
	}
	if good, err = os.OpenFile(goodfile, os.O_RDWR, 0644); err != nil {
		panic(err)
	}
	if bad, err = os.OpenFile(badfile, os.O_RDWR, 0644); err != nil {
		panic(err)
	}
	defer func() {
		good.Close()
		raw.Close()
		bad.Close()
	}()
	reader := csv.NewReader(raw)
	_, err = reader.Read()
	Len := 0
	for {
		_, err = reader.Read()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		Len++
	}
	wg := sync.WaitGroup{}
	goodwriter := csv.NewWriter(good)
	badwriter := csv.NewWriter(bad)
	notification := time.NewTicker(time.Second * 5)
	done := make(chan any)
	task := make(chan any, 100)
	progress := make(chan any, 5000)
	go func() {
		for {
			select {
			case <-done:
				{
					return
				}
			case <-notification.C:
				{
					fmt.Printf("执行进度 %d/%d %d%% \n", len(progress), Len, 100*len(progress)/Len)
				}
			}
		}
	}()
	client := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Second * 5,
		},
	}
	raw.Seek(0, 0)
	safe := int64(0)
	var line []string
	_, _ = reader.Read()
	var GLock, BLock sync.RWMutex
	for {
		line, err = reader.Read()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		wg.Add(1)
		go func(raw []string, url string) {
			task <- nil
			defer wg.Done()
			defer func() { <-task }()
			resp, err := client.Get(url)
			progress <- nil
			if err != nil || resp.StatusCode != http.StatusOK {
				BLock.Lock()
				defer BLock.Unlock()
				badwriter.Write(raw)
				badwriter.Flush()
				return
			}
			GLock.Lock()
			defer GLock.Unlock()
			goodwriter.Write(raw)
			goodwriter.Flush()
			defer resp.Body.Close()
			atomic.AddInt64(&safe, 1)
		}(line, line[4])
	}
	wg.Wait()
	done <- nil
	fmt.Printf("Good URL ：%d \n", safe)
}
func main() {
	start := time.Now()
	fmt.Println("开始执行任务")
	Task()
	fmt.Printf("共计用时：%v \n", time.Now().Sub(start))
}
