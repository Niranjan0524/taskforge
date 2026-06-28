package config

import (
	"net/url"
	"os"
	"strings"
)

func OriginAllowed(origin string) bool {
	origin = normalizeOrigin(origin)
	if origin == "" {
		return true
	}

	allowedOrigins := os.Getenv("ORIGIN_URL")
	if allowedOrigins == "" {
		return true
	}

	for _, allowed := range strings.Split(allowedOrigins, ",") {
		allowed = normalizeOrigin(allowed)
		if allowed == "*" || allowed == origin {
			return true
		}
	}

	host := originHost(origin)
	return host == "localhost" ||
		strings.HasSuffix(host, ".localhost") ||
		strings.HasSuffix(host, ".vercel.app")
}

func normalizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	origin = strings.TrimRight(origin, "/")

	return origin
}

func originHost(origin string) string {
	parsed, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	return strings.ToLower(host)
}
