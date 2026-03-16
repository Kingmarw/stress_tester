package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	requests := 1000000 // 1 مليون طلب
	concurrency := 1000 // عدد الـ goroutines اللي هتشتغل في نفس الوقت

	start := time.Now()
	fmt.Println("Starting stress test...")

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency) // للتحكم في الـ concurrency

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{} // حجز goroutine
		go func() {
			defer wg.Done()
			http.Get("https://ts-oracle.co/erp_aqua/") // رابط السيرفر بتاعك
			<-sem // تحرير goroutine
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("Total Requests:", requests)
	fmt.Println("Time:", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(requests)/duration.Seconds())
}