package main

import "os"

type config struct {
	data map[string]bool
}

func (cm *config) add(site string) {
	cm.data[site] = true
}

func (cm *config) exists(site string) bool {
	return cm.data[site]
}

func (cm *config) delete(site string) {
	cm.data[site] = false
}

func (cm *config) save(filename string) error {
	var data string
	for key, yes := range cm.data {
		if yes {
			data += key + "\n"
		}
	}
	return os.WriteFile(filename, []byte(data), 0644)
}
