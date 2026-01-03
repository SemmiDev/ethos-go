package httputil

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP tries to find the real client IP from httputil headers or connection info
func GetClientIP(r *http.Request) string {
	// Check common proxy headers (in order)
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
		"X-Client-IP",      // Some proxies
	}

	for _, h := range headers {
		ips := strings.SplitSeq(r.Header.Get(h), ",")
		for ip := range ips {
			ip = strings.TrimSpace(ip)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetUserAgent safely gets the User-Agent header
func GetUserAgent(r *http.Request) string {
	ua := strings.TrimSpace(r.Header.Get("User-Agent"))
	if ua == "" {
		return "Unknown"
	}
	return ua
}
