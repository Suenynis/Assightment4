package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, j := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go j(in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range in {
		data := fmt.Sprintf("%v", val)
		md5 := DataSignerMd5(data)
		crc32 := DataSignerCrc32(data)
		out <- fmt.Sprintf("%s~%s", md5, crc32)
	}
}

func MultiHash(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range in {
		data := fmt.Sprintf("%v", val)
		var crc32results []string
		for i := 0; i < 6; i++ {
			crc32 := DataSignerCrc32(strconv.Itoa(i) + data)
			crc32results = append(crc32results, crc32)
		}
		out <- strings.Join(crc32results, "")
	}
}

func CombineResults(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var results []string
	for val := range in {
		data := fmt.Sprintf("%v", val)
		results = append(results, data)
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}
