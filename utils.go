package main

import (
	"bufio"
	"os"
	"strings"
)

func parseSocialsToMap(hostFile string, sm map[string][]string) error {
	file, err := os.Open(hostFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var dName string
	var dUrls []string
	var SaveToMap = func() {
		if dName != "" && len(dUrls) > 0 {
			sm[dName] = dUrls
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "# [") && strings.HasSuffix(line, "]") {
			SaveToMap() // save current holding key
			dName = line[3 : len(line)-1]
			dUrls = []string{}
		} else if strings.HasPrefix(line, "#") {
			continue
		} else {
			fields := strings.Fields(line)
			if len(fields) > 1 && fields[0] == "0.0.0.0" {
				dUrls = append(dUrls, fields[1])
			}
		}
	}
	SaveToMap()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func parseConfigToMap(filename string, cm map[string]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if state := cm[key]; !state {
			cm[key] = true
		}
	}
	return nil
}
