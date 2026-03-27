import sys
import requests

def run_audit(url):
    try:
        # بنقلد المتصفح عشان Vercel ميعملش Block
        headers = {'User-Agent': 'Mozilla/5.0'}
        report = []
        
        # 1. فحص ملفات المطورين (Common on Vercel/Next.js)
        check_paths = ['/robots.txt', '/sitemap.xml', '/.env']
        for path in check_paths:
            r = requests.get(url + path, headers=headers, timeout=5)
            if r.status_code == 200:
                report.append(f"Found: {path}")

        if not report:
            return "Scan complete: No sensitive files exposed in common paths."
        return " | ".join(report)

    except Exception as e:
        return f"Logic Error: {str(e)}"

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 deep_scan.py <url>")
        sys.exit(1)
    
    target = sys.argv[1]
    # تنظيف الرابط من أي فراغات خفية
    target = target.strip()
    print(run_audit(target))
