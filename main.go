package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TestResults struct {
	mu           sync.Mutex
	totalReqs    int
	success      int
	fail         int
	netErrors    int
	statusCodes  map[int]int
	totalLatency time.Duration
	minLatency   time.Duration
	maxLatency   time.Duration
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
}

func getRandomAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func getRandomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}


func spoofBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", getRandomAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("X-Forwarded-For", getRandomIP()) 
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printLogo() {
	logo := `	 _       _________ _        _______  _______  _______            _________ _        _______ 
	| \    /\\__   __/( (    /|(  ____ \(       )(  ____ )|\     /|  \__   __/( (    /|(  ____ \
	|  \  / /   ) (   |  \  ( || (    \/| () () || (    )|| )   ( |     ) (   |  \  ( || (    \/
	|  (_/ /    | |   |   \ | || |      | || || || (____)|| | _ | |     | |   |   \ | || (_____ 
	|   _ (     | |   | (\ \) || | ____ | |(_)| ||     __)| |( )| |     | |   | (\ \) |(_____  )
	|  ( \ \    | |   | | \   || | \_  )| |   | || (\ (   | || || |     | |   | | \   |      ) |
	|  /  \ \___) (___| )  \  || (___) || )   ( || ) \ \__| () () |  ___) (___| )  \  |/\____) |
	|_/    \/\_______/|/    )_)(_______)|/     \||/   \__/(_______)  \_______/|/    )_)\_______)
                                                                                                                                                                         
    [ (Demo) V1 PHANTOM EDITION - ADVANCED OSINT, FUZZING & CYBER INTELLIGENCE ]
=========================================================================================`
	fmt.Println(logo)
}

func runDeepRecon(targetURL string, client *http.Client) {
	fmt.Println("\n🧠 [Phase 1] Initializing Deep AI Reconnaissance & Spoofing...")

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("   ❌ Invalid URL format.")
		return
	}
	hostname := parsedURL.Hostname()

	fmt.Println("\n   🌍 Network & Topology:")
	ips, err := net.LookupIP(hostname)
	if err == nil && len(ips) > 0 {
		fmt.Printf("      🔹 Target IP(s): ")
		for _, ip := range ips {
			fmt.Printf("%s ", ip.String())
		}
		fmt.Println()
	} else {
		fmt.Println("      🔹 Target IP: Hidden or Unresolvable")
	}

	req, _ := http.NewRequest("GET", targetURL, nil)
	spoofBrowserHeaders(req) // تفعيل وضع التمويه
	
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ❌ Target unreachable: %v\n", err)
		return
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	bodyString := string(bodyBytes)
	bodyLower := strings.ToLower(bodyString)

	cookies := ""
	for _, c := range res.Cookies() {
		cookies += c.Name + "=" + c.Value + ";"
	}

	serverHeader := res.Header.Get("Server")

	// WAF & Deception Detection (الرادار الذكي)
	fmt.Println("\n   🛡️  WAF & Anti-Bot Shield:")
	wafFound := false
	if strings.Contains(strings.ToLower(serverHeader), "cloudflare") || strings.Contains(cookies, "__cf") {
		fmt.Println("      ⚠️  Protected by: Cloudflare (Advanced WAF)")
		wafFound = true
	}
	if strings.Contains(serverHeader, "Akamai") {
		fmt.Println("      ⚠️  Protected by: Akamai CDN/WAF")
		wafFound = true
	}
	if strings.Contains(bodyLower, "aes.js") || strings.Contains(cookies, "__test") {
		fmt.Println("      🚨 DECEPTION : AES Anti-Bot Shield (InfinityFree/iFastNet)")
		fmt.Println("      ⚠️  WARNING  : Source code is masked. Scan accuracy may drop.")
		wafFound = true
	}
	if strings.Contains(serverHeader, "Sucuri") {
		fmt.Println("      ⚠️  Protected by: Sucuri Cloudproxy")
		wafFound = true
	}
	if !wafFound {
		fmt.Println("      ✅ No standard WAF detected (Direct access possible)")
	}

	// Scraper
	fmt.Println("\n   🕵️  Data Extraction:")
	titleRegex := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	if match := titleRegex.FindStringSubmatch(bodyString); len(match) > 1 {
		fmt.Printf("      📄 Title  : %s\n", strings.TrimSpace(match[1]))
	}
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(bodyString, -1)
	if len(emails) > 0 {
		uniqueEmails := make(map[string]bool)
		fmt.Printf("      📧 Emails : ")
		for _, e := range emails {
			if !uniqueEmails[e] {
				fmt.Printf("%s ", e)
				uniqueEmails[e] = true
			}
		}
		fmt.Println()
	}

	// Backend & Core
	fmt.Println("\n   ⚙️  Backend & Core Technologies:")
	if serverHeader != "" {
		fmt.Printf("      🔹 Web Server : %s\n", serverHeader)
	}
	if strings.Contains(res.Header.Get("X-Powered-By"), "PHP") || strings.Contains(cookies, "PHPSESSID") || strings.Contains(bodyLower, ".php") {
		fmt.Println("      🔹 Language   : PHP")
	}
	if strings.Contains(cookies, "ASP.NET") || res.Header.Get("X-AspNet-Version") != "" {
		fmt.Println("      🔹 Language   : C# / ASP.NET")
	}
	if strings.Contains(cookies, "connect.sid") || strings.Contains(res.Header.Get("X-Powered-By"), "Express") {
		fmt.Println("      🔹 Language   : Node.js (Express)")
	}
	if strings.Contains(cookies, "JSESSIONID") {
		fmt.Println("      🔹 Language   : Java")
	}

	// Advanced Frontend Detection
	fmt.Println("\n   🎨 Frontend Architecture & Frameworks:")
	uiDetected := false
	if strings.Contains(bodyLower, "wp-content") {
		fmt.Println("      🔸 CMS        : WordPress")
		uiDetected = true
	}
	if regexp.MustCompile(`id="__next"|/_next/static`).MatchString(bodyLower) {
		fmt.Println("      🔸 Framework  : Next.js (React)")
		uiDetected = true
	} else if regexp.MustCompile(`data-reactroot|react-dom`).MatchString(bodyLower) {
		fmt.Println("      🔸 Framework  : React.js")
		uiDetected = true
	}
	if regexp.MustCompile(`id="__nuxt"|/_nuxt/`).MatchString(bodyLower) {
		fmt.Println("      🔸 Framework  : Nuxt.js (Vue)")
		uiDetected = true
	} else if regexp.MustCompile(`data-v-[a-z0-9]+`).MatchString(bodyLower) {
		fmt.Println("      🔸 Framework  : Vue.js")
		uiDetected = true
	}
	// التقنيات الجديدة
	if strings.Contains(bodyLower, "_astro") {
		fmt.Println("      🔸 Framework  : Astro.build (Modern SSG)")
		uiDetected = true
	}
	if strings.Contains(bodyLower, "svelte-") {
		fmt.Println("      🔸 Framework  : Svelte / SvelteKit")
		uiDetected = true
	}
	if strings.Contains(bodyLower, "vite/client") || strings.Contains(bodyLower, "@vite") {
		fmt.Println("      🔸 Bundler    : Vite.js")
		uiDetected = true
	}

	if regexp.MustCompile(`class="[^"]*\b(flex|grid|p-[0-9]|text-center|bg-[a-z]+-[0-9]+)\b[^"]*"`).MatchString(bodyLower) {
		fmt.Println("      🔸 UI Styling : Tailwind CSS")
		uiDetected = true
	}
	if regexp.MustCompile(`class="[^"]*\b(container|row|col-|btn-|navbar)\b[^"]*"`).MatchString(bodyLower) {
		fmt.Println("      🔸 UI Styling : Bootstrap")
		uiDetected = true
	}
	if !uiDetected {
		fmt.Println("      🔸 Architecture: Native HTML/CSS/JS (Or obfuscated)")
	}

	fmt.Println("\n   🔒 Security Headers Audit:")
	headers := map[string]string{
		"Strict-Transport-Security": "HSTS",
		"X-Frame-Options":           "Clickjacking",
		"Content-Security-Policy":   "XSS Protection",
	}
	for h, desc := range headers {
		if res.Header.Get(h) != "" {
			fmt.Printf("      ✅ %-25s : Secured\n", h)
		} else {
			fmt.Printf("      ❌ %-25s : VULNERABLE (%s)\n", h, desc)
		}
	}
    fmt.Println("\n🐍 [Extra] Running Python Deep Scan Plugin...")
    pythonResult := runPythonPlugin(targetURL)
    fmt.Printf("    📑 Python Intelligence Report: %s\n", pythonResult)
	fmt.Println("=========================================================================================")
}

func runDeepFuzzer(targetURL string, client *http.Client) {
	fmt.Println("\n💫 [Phase 2] Executing Deep Directory Fuzzing & Anomaly Detection...")
	baseURL := strings.TrimRight(targetURL, "/")
	
	hiddenPaths := []string{
		"/.env", "/.git/config", "/wp-config.php.bak", "/config.json",
		"/admin", "/dashboard", "/api/v1/users", "/swagger/v1/swagger.json", 
		"/phpinfo.php", "/robots.txt", "/sitemap.xml", "/.DS_Store",
	}

	var wg sync.WaitGroup
	for _, path := range hiddenPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			req, _ := http.NewRequest("GET", baseURL+p, nil)
			spoofBrowserHeaders(req) 
			
			res, err := client.Do(req)
			if err == nil {
				// Honeypot heuristic: If an obscure file returns 200 but content length is huge, it might be a honeypot catching bots.
				if res.StatusCode == 200 {
					fmt.Printf("   🚨 FATAL EXPOSURE : %-25s (HTTP 200 OK)\n", p)
				} else if res.StatusCode == 403 {
					fmt.Printf("   🔒 RESTRICTED     : %-25s (HTTP 403 - Exists)\n", p)
				} else if res.StatusCode == 401 {
					fmt.Printf("   🔑 AUTH REQUIRED  : %-25s (HTTP 401 - Valid Target)\n", p)
				}
				res.Body.Close()
			}
		}(path)
		time.Sleep(150 * time.Millisecond) // Evasion delay
	}
	wg.Wait()
	fmt.Println("=========================================================================================")
}

func runLoadTest(url string, requests int, concurrency int, client *http.Client) {
	fmt.Printf("\n⚔️ [Phase 3] Initiating Phantom Stress Attack on: %s\n", url)
	
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	results := &TestResults{statusCodes: make(map[int]int), minLatency: time.Hour}
	startTime := time.Now()

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			reqStart := time.Now()

			req, _ := http.NewRequest("GET", url, nil)
			spoofBrowserHeaders(req) // هجوم بـ IPs وهمية

			res, err := client.Do(req)
			reqDuration := time.Since(reqStart)

			results.mu.Lock()
			results.totalReqs++
			if reqDuration < results.minLatency { results.minLatency = reqDuration }
			if reqDuration > results.maxLatency { results.maxLatency = reqDuration }
			results.totalLatency += reqDuration

			if err != nil {
				results.netErrors++
			} else {
				results.statusCodes[res.StatusCode]++
				if res.StatusCode >= 200 && res.StatusCode < 300 {
					results.success++
				} else {
					results.fail++
				}
				res.Body.Close()
			}
			results.mu.Unlock()
			<-sem
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)

	fmt.Println("\n================================ 📈 MISSION REPORT ================================")
	fmt.Printf("   🚀 Raw Speed      : %.2f req/sec\n", float64(results.totalReqs)/totalTime.Seconds())
	avgLatency := time.Duration(0)
	if results.totalReqs-results.netErrors > 0 {
		avgLatency = time.Duration(int64(results.totalLatency) / int64(results.totalReqs-results.netErrors))
	}
	fmt.Printf("   ⏱️  Server Ping    : %v (Fastest: %v | Slowest: %v)\n", avgLatency, results.minLatency, results.maxLatency)
	fmt.Println("-----------------------------------------------------------------------------------------")
	fmt.Printf("   ✅ Hits: %d  |  ❌ Blocks/Fails: %d  |  ⚠️ Dropped: %d\n", results.success, results.fail, results.netErrors)
	
	if len(results.statusCodes) > 0 {
		fmt.Println("\n   📌 Server Response Breakdown:")
		for code, count := range results.statusCodes {
			fmt.Printf("      [ HTTP %d ] -> Processed %d times\n", code, count)
		}
	}
	fmt.Println("=========================================================================================")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printLogo()

		fmt.Print("\n🔗 Enter Target URL (e.g., example.com): ")
		urlInput, _ := reader.ReadString('\n')
		urlInput = strings.TrimSpace(urlInput)
		
		// Auto-fix URL 
		urlInput = strings.ReplaceAll(urlInput, "https://http//", "http://")
		if urlInput != "" && !strings.HasPrefix(urlInput, "http://") && !strings.HasPrefix(urlInput, "https://") {
			urlInput = "https://" + urlInput
		}

		if urlInput == "" || urlInput == "https://" {
			fmt.Println("❌ Error: Valid URL is required.")
			time.Sleep(2 * time.Second)
			continue
		}

		client := &http.Client{Timeout: 10 * time.Second}

		fmt.Println("\n🔥 Select Operation Mode:")
		fmt.Println("   1. 🧠 Deep Recon (WAF, OSINT, Tech Core & Deception Check)")
		fmt.Println("   2. 💫 Phantom Fuzzer (Hunt for Hidden Routes)")
		fmt.Println("   3. ⚔️ Load Tester (Spoofed IPs Stress Attack)")
		fmt.Println("   4. 🚀 APEX MODE (Full Automated Audit)")
		fmt.Print("\n👉 Enter choice (1-4): ")
		
		choiceStr, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(choiceStr)

		switch choice {
		case "1":
			runDeepRecon(urlInput, client)
		case "2":
			runDeepFuzzer(urlInput, client)
		case "3":
			fmt.Print("   ⚡ Total Requests (Default 100): ")
			reqStr, _ := reader.ReadString('\n')
			reqs, _ := strconv.Atoi(strings.TrimSpace(reqStr))
			if reqs <= 0 { reqs = 100 }

			fmt.Print("   🚀 Concurrency (Default 10): ")
			conStr, _ := reader.ReadString('\n')
			conc, _ := strconv.Atoi(strings.TrimSpace(conStr))
			if conc <= 0 { conc = 10 }
			
			runLoadTest(urlInput, reqs, conc, client)
		case "4":
			runDeepRecon(urlInput, client)
			runDeepFuzzer(urlInput, client)
			runLoadTest(urlInput, 200, 20, client)
		default:
			fmt.Println("❌ Invalid choice. Returning to menu...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("\n🔄 Press [Enter] to scan another target, or type 'q' to Quit: ")
		exitCmd, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(exitCmd)) == "q" {
			fmt.Println("\n👋 Exiting Kingmarw Phantom Framework... Stay Sharp!")
			break
		}
	}
}
func runPythonPlugin(target string) string {
	// تنفيذ أمر: python3 deep_scan.py target
	cmd := exec.Command("python3", "deep_scan.py", target)
	
	// سحب النتيجة اللي بايثون طبعتها (Stdout)
	out, err := cmd.Output()
	if err != nil {
		return "Python Plugin Error: " + err.Error()
	}
	
	return strings.TrimSpace(string(out))
}
