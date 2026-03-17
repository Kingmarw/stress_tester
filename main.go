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

	// متغيرات لعد حالة الطلبات
	var successCount int64
	var failCount int64
	var errorCount int64
	var printedErrors int64 // متغير جديد عشان نطبع أول 5 أخطاء بس

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{} // حجز goroutine
		
		go func(url string) {
			defer wg.Done()
			
			// تجهيز الطلب وإضافة User-Agent وهمي عشان نخدع حماية السيرفرات
			req, err := http.NewRequest("GET", url, nil)
			if err == nil {
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
			}
			
			// تنفيذ الطلب
			res, err := http.DefaultClient.Do(req)
			
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				<-sem // تحرير goroutine
				return
			}
			
			// فحص النتيجة
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				atomic.AddInt64(&successCount, 1) // ناجح
			} else {
				atomic.AddInt64(&failCount, 1) // فاشل
				
				// طباعة أول 5 أخطاء بس عشان نعرف السبب
				if atomic.AddInt64(&printedErrors, 1) <= 5 {
					fmt.Printf("[!] Request failed - Status Code: %d\n", res.StatusCode)
				}
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
