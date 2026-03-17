package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// استقبال الرابط من المستخدم
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the target URL (e.g., https://example.com/): ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)

	// التأكد إن الرابط مش فاضي
	if targetURL == "" {
		fmt.Println("Error: URL cannot be empty.")
		return
	}

	requests := 1000000 // 1 مليون طلب
	concurrency := 1000 // عدد الـ goroutines اللي هتشتغل في نفس الوقت

	start := time.Now()
	fmt.Printf("\nStarting stress test on %s...\n", targetURL)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency) // للتحكم في الـ concurrency

	// متغيرات لعد حالة الطلبات (استخدمنا int64 عشان نستخدمها مع atomic)
	var successCount int64
	var failCount int64
	var errorCount int64

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{} // حجز goroutine
		
		go func(url string) {
			defer wg.Done()
			
			// عمل الطلب
			res, err := http.Get(url)
			
			// لو حصل خطأ في الاتصال نفسه (زي إن السيرفر وقع أو النت فصل)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				<-sem // تحرير goroutine
				return
			}
			
			// لو الطلب وصل، نتأكد من الـ Status Code
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				atomic.AddInt64(&successCount, 1) // ناجح (مثلا 200 OK)
			} else {
				atomic.AddInt64(&failCount, 1) // فاشل (مثلا 503 أو 429)
			}

			// قفل الاتصال
			res.Body.Close()
			
			<-sem // تحرير goroutine
		}(targetURL)
	}

	wg.Wait()
	duration := time.Since(start)

	// طباعة التقرير النهائي
	fmt.Println("\n--- Test Completed ---")
	fmt.Println("Target URL:", targetURL)
	fmt.Println("Total Requests:", requests)
	fmt.Println("----------------------")
	fmt.Printf("Successful (2xx): %d\n", successCount)
	fmt.Printf("Failed (Non-2xx): %d\n", failCount)
	fmt.Printf("Network Errors:   %d\n", errorCount)
	fmt.Println("----------------------")
	fmt.Println("Time:", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(requests)/duration.Seconds())
}
