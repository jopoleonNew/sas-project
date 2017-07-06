package main

import (
	"fmt"
	"sync"
)

const xthreads = 3 // Total number of threads to use, excluding the main() thread

func doSomething(a int) {
	fmt.Println("My job is", a)
	return
}

func main() {
	var ch = make(chan int, 3) // This number 50 can be anything as long as it's larger than xthreads
	var wg sync.WaitGroup

	// This starts xthreads number of goroutines that wait for something to do
	wg.Add(xthreads)
	for i := 0; i < xthreads; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				doSomething(a) // do the thing
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for i := 0; i < 50; i++ {
		ch <- i // add i to the queue
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
}

//
//type idProcessor func(id uint)
//
//func SpawnStuff(limit uint, proc idProcessor) chan<- uint {
//	ch := make(chan uint)
//	for i := uint(0); i < limit; i++ {
//		go func() {
//			for {
//				id, ok := <-ch
//				if !ok {
//					return
//				}
//				proc(id)
//			}
//		}()
//	}
//	return ch
//}
//
//func main() {
//	runtime.GOMAXPROCS(4)
//	var wg sync.WaitGroup //this is just for the demo, otherwise main will return
//	fn := func(id uint) {
//		fmt.Print(" ", id)
//		wg.Done()
//	}
//	wg.Add(1000)
//	ch := SpawnStuff(10, fn)
//	for i := uint(0); i < 1000; i++ {
//		ch <- i
//	}
//	close(ch) //should do this to make all the goroutines exit gracefully
//	wg.Wait()
//}

//func main() {
//	maxGoroutines := 3
//	guard := make(chan struct{}, maxGoroutines)
//
//	for i := 0; i < 30; i++ {
//		guard <- struct{}{} // would block if guard channel is already filled
//		go func(i int) {
//			worker(i)
//			<-guard
//		}(i)
//	}
//}
//
//func worker(i int) { fmt.Println("doing work on", i) }

//
//func main() {
//	urls := []string{
//		"http://www.reddit.com/r/aww.json",
//		"http://www.reddit.com/r/funny.json",
//		"http://www.reddit.com/r/programming.json",
//	}
//	jsonResponses := make(chan string)
//
//	var wg sync.WaitGroup
//
//	wg.Add(len(urls))
//
//	for _, url := range urls {
//		go func(url string) {
//			defer wg.Done()
//			res, err := http.Get(url)
//			if err != nil {
//				log.Fatal(err)
//			} else {
//				defer res.Body.Close()
//				body, err := ioutil.ReadAll(res.Body)
//				if err != nil {
//					log.Fatal(err)
//				} else {
//					jsonResponses <- string(body)
//				}
//			}
//		}(url)
//	}
//
//	go func() {
//		for response := range jsonResponses {
//			fmt.Println(response)
//		}
//	}()
//
//	wg.Wait()
//}
