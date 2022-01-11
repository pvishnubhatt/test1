package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var msg int = 1

func main() {
	producerConsumer()
}

var ops uint64

func producerConsumer() {
	data := make(chan uint64)

	go cproducer(data)
	go cconsumer(data)
	select {}
}

func cproducer(data chan uint64) {
	for {
		atomic.AddUint64(&ops, 1)
		data <- ops
		fmt.Printf("Produced %d ", ops)
	}
}

func cconsumer(data chan uint64) {
	for {
		xx := <-data
		fmt.Printf(", Consumed %d\n", xx)
	}
}

func testUnbufferedChannel() {
	fmt.Println("In testUnbufferedChannel")
	data := make(chan string)
	go func() {
		data <- "Ping"
	}()
	for {
		go pingpong(data)
	}
}

func pingpong(data chan string) {
	stat := <-data
	fmt.Printf("%s\n", stat)
	dstat := stat
	if stat == "Ping" {
		dstat = "Pong"
	} else {
		dstat = "Ping"
	}
	data <- dstat
}

func testWaitGroup() {
	fmt.Println("In testWaitGroup")
	var wg sync.WaitGroup
	go producer(&wg)
	go consumer(&wg)
	time.Sleep(time.Microsecond)
	wg.Wait()
	fmt.Println("Done testWaitGroup")
}

func producer(wg *sync.WaitGroup) {
	for {
		wg.Add(1)
		//fmt.Printf("Produced %d\n", *msg)
		msg += 1
	}
}

func consumer(wg *sync.WaitGroup) {
	for {
		wg.Done()
		fmt.Printf("Consumed %d\n", msg)
	}
}

func test1(wg *sync.WaitGroup) {
	fmt.Println("In test1")
	wg.Done()
}
