package http

import (
	"fmt"
	"strings"
)

var knownContent map[string]string

func init() {
	knownContent = map[string]string{
		".js":     "application/javascript",
		".css":    "text/css",
		".ico":    "image/x-icon",
		".gif":    "image/gif",
		"js.map":  "application/octet-stream",
		"css.map": "application/octet-stream",
		".woff":   "font/woff",
		".woff2":  "font/woff2",
		".ttf":    "font/ttf",
	}
}

func detectType(filename string) (string, error) {
	for suffix, contentType := range knownContent {
		if strings.HasSuffix(filename, suffix) {
			return contentType, nil
		}
	}

	return "", fmt.Errorf("unknown type %v", filename)
}
