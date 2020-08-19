package impl

import (
	"fmt"
	"strings"
)

type Content struct {
	mytype map[string]string
}

func NewContent() *Content {
	c := &Content{}
	c.mytype = map[string]string{
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
	return c
}

func (c *Content) detectType(filename string) (string, error) {
	for suffix, contentType := range c.mytype {
		if strings.HasSuffix(filename, suffix) {
			return contentType, nil
		}
	}

	return "", fmt.Errorf("unknown type %v", filename)
}
