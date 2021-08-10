package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(myJob ...job) {
	channel := make(chan interface{}, 100)
	wg := &sync.WaitGroup{}

	for _, j := range myJob {
		wg.Add(1)
		channel = returnChannel(j, channel, wg)
	}
	wg.Wait()
}

func returnChannel(j job, c chan interface{}, wg *sync.WaitGroup) chan interface{} {
	channel := make(chan interface{}, 100)
	go func(j job) {
		defer wg.Done()
		defer close(channel)
		j(c, channel)
	}(j)
	return channel
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for data := range in {
		str := fmt.Sprintf("%v", data)
		wg.Add(1)
		go func(str string, quotaMu *sync.Mutex) {
			defer wg.Done()

			resMd5 := make(chan string)
			go func() {
				quotaMu.Lock()
				resMd5 <- DataSignerMd5(str)
				quotaMu.Unlock()

			}()

			result := make(chan string)
			go func() {
				result <- DataSignerCrc32(str)
			}()
			str2 := DataSignerCrc32(<-resMd5)
			str1 := <-result
			out <- str1 + "~" + str2
		}(str, mu)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		str := fmt.Sprintf("%v", data)
		mu := &sync.Mutex{}
		result := make([]string, 6)
		wg.Add(1)
		go func(data string, wg *sync.WaitGroup) {
			defer wg.Done()
			wgInside := &sync.WaitGroup{}
			for i := 0; i < 6; i++ {
				th := fmt.Sprintf("%v", i)
				wgInside.Add(1)
				go func(num int, wgInside *sync.WaitGroup) {
					defer wgInside.Done()
					res := DataSignerCrc32(th + data)
					mu.Lock()
					result[num] = res
					mu.Unlock()
				}(i, wgInside)
			}
			wgInside.Wait()
			out <- strings.Join(result, "")
		}(str, wg)
	}
	wg.Wait()
	out <- 1
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0)
LOOP:
	for {
		select {
		case val := <-in:
			if val == 1 {
				sort.Strings(result)
				out <- strings.Join(result, "_")
				break LOOP
			} else {
				result = append(result, fmt.Sprintf("%v", val))
			}
		}
	}
}
