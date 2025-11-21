package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func populateSocialMap(filename string, sm GhostMap) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var domain string
	var urls []string
	var saveMap = func() {
		if domain != "" && len(urls) > 0 {
			sm.data[domain] = urls
		}
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "# [") && strings.HasSuffix(line, "]") {
			saveMap()
			domain = line[3 : len(line)-1]
			urls = []string{}
		} else if strings.HasPrefix(line, "#") {
			continue
		} else {
			fields := strings.Fields(line)
			if len(fields) > 1 && fields[0] == "0.0.0.0" {
				urls = append(urls, fields[1])
			} else {
				urls = append(urls, line)
			}
		}
	}
	saveMap()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func populateConfigMap(filename string, cm GhostMap, gm GhostMap) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if len(domain) == 0 {
			continue
		}
		if _, ok := gm.data[domain]; !ok {
			return fmt.Errorf("config: %v, no such Host", domain)
		}
		cm.data[domain] = gm.data[domain]
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func printListView(sm, cm GhostMap) string {
	var builder strings.Builder
	for domain := range sm.data {
		if _, ok := cm.data[domain]; !ok {
			builder.WriteString(fmt.Sprintf("%s%s%s\n", Red, domain, Reset))
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", domain))
		}
	}
	return strings.TrimSpace(builder.String())
}
