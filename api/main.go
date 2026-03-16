package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// استقبال الرابط من المستخدم
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the target URL (e.g., https://example.com/): ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL) // تنظيف الرابط من أي مسافات

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

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{} // حجز goroutine
		
		// تمرير الرابط كمتغير للـ goroutine لضمان استقرار الأداء
		go func(url string) {
			defer wg.Done()
			
			// عمل الطلب وقفل الاتصال فوراً لتجنب استهلاك كل بورتات الجهاز (Socket Exhaustion)
			res, err := http.Get(url)
			if err == nil {
				res.Body.Close() 
			}
			
			<-sem // تحرير goroutine
		}(targetURL)
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("\n--- Test Completed ---")
	fmt.Println("Target URL:", targetURL)
	fmt.Println("Total Requests:", requests)
	fmt.Println("Time:", duration)
	fmt.Printf("Requests/sec: %.2f\n", float64(requests)/duration.Seconds())
}
