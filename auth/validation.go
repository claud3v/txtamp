package auth

import (
	"net/url"
	"strings"
)

func IsNotBlank(username string) bool {
	return strings.TrimSpace(username) != ""
}

func IsServerUrlValid(urlInput string) bool {
	trimmed := strings.TrimSpace(urlInput)
	if trimmed == "" {
		return false
	}

	u, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
