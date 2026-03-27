- Nikto v2.6.0
---------------------------------------------------------------------------
+ Your Nikto installation is out of date.
+ Target IP:          64.29.17.67
+ Target Hostname:    kingmarw.vercel.app
+ Target Port:        443
---------------------------------------------------------------------------
+ SSL Info:           Subject:  /CN=*.vercel.app
                      CN:       *.vercel.app
                      SAN:      *.vercel.app, vercel.app
                      Ciphers:  TLS_AES_128_GCM_SHA256
                      Issuer:   /C=US/O=Google Trust Services/CN=WR1
+ Platform:           Unknown
+ Start Time:         2026-03-26 21:31:14 (GMT-4)
---------------------------------------------------------------------------
+ Server: Vercel
+ Multiple IPs found: 64.29.17.67, 216.198.79.67
+ [999986] /: Retrieved access-control-allow-origin header: *.
+ [999979] /: IP address found in the 'x-vercel-id' header. The IP is "a1::". See: https://portswigger.net/kb/issues/00600300_private-ip-addresses-disclosed
+ [999100] /: Uncommon header(s) 'x-vercel-id' found, with contents: fra1::vbhcp-1774575074646-7b25c2acf017.
+ [999100] /: Uncommon header(s) 'content-disposition' found, with contents: inline.
+ [999100] /: Uncommon header(s) 'x-vercel-cache' found, with contents: HIT.
+ [999100] /server-status: Uncommon header(s) 'x-vercel-error' found, with contents: NOT_FOUND.
+ [999100] /cgi-bin/: Uncommon header(s) 'x-vercel-mitigated' found, with contents: deny.
+ No CGI Directories found (use '-C all' to force check all possible dirs). CGI tests skipped.
+ [999996] /robots.txt: contains 1 entry which should be manually viewed. See: https://developer.mozilla.org/en-US/docs/Glossary/Robots.txt
+ ERROR: *** Error limit (20) reached for host, giving up. Last error: opening stream: ssl connect failed. ***
+ ERROR: *** Consider using mitmproxy to avoid TLS fingerprinting. ***
- STATUS: Completed 107 requests (~3% complete, ~31.7 minutes left): currently in plugin 'Site Files'
- STATUS: Running average: Not enough data.
+ Scan terminated: 20 errors and 8 items reported on the remote host
+ End Time:           2026-03-26 21:32:10 (GMT-4) (56 seconds)
---------------------------------------------------------------------------
+ 1 host(s) tested
