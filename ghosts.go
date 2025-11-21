package main

import (
	"fmt"
	"os"
	"strings"
)

type GhostMap struct {
	data map[string][]string
}

func (gm *GhostMap) String() string {
	var builder strings.Builder
	for domain, urls := range gm.data {
		builder.WriteString(fmt.Sprintf("# [%s]\n", domain))
		if len(urls) > 0 {
			for _, u := range urls {
				builder.WriteString(fmt.Sprintf("0.0.0.0 %s\n", u))
			}
		}
		builder.WriteString("\n")
	}
	result := builder.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

func (gm *GhostMap) IsExists(key string) bool {
	_, ok := gm.data[key]
	return ok
}

func (gm *GhostMap) SaveToFile(filename string) error {
	return os.WriteFile(filename, []byte(gm.String()), 0644)
}

func (gm *GhostMap) List() string {
	var builder strings.Builder
	for domain := range gm.data {
		builder.WriteString(fmt.Sprintf("%v\n", domain))
	}
	return strings.TrimSpace(builder.String())
}
