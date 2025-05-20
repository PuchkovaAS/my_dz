package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	nums := make(chan int)
	results := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)
	go createRandSlice(nums, &wg)
	go powSlice(nums, results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Print(res, " ")
	}
}

func createRandSlice(ch chan int, wg *sync.WaitGroup) {
	for range 10 {
		ch <- rand.Intn(100)
	}
	wg.Done()
	close(ch)
}

func powSlice(nums <-chan int, results chan<- int, wg *sync.WaitGroup) {
	for num := range nums {
		results <- num * num
	}
	wg.Done()
}
