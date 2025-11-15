package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

// const ETC_HOSTS = "/etc/hosts"
const ETC_HOSTS = "internal/etc.hosts"

var defaultGhosts = "internal/default.hosts"
var hostFile = "internal/social.hosts"
var configFile = "internal/config.hosts"

var socialsMap = make(map[string][]string)
var configMap = config{
	data: make(map[string]bool),
}

var list = flag.Bool("list", false, "list running ghost sites")
var add = flag.String("add", "", "add a ghost site")
var delete = flag.String("del", "", "remove a ghost site")
var listAll = flag.Bool("la", false, "list all ghost sites")

func main() {
	switch {
	case len(os.Args) == 1:
		flag.Usage()
	case *list:
		for key, yes := range configMap.data {
			if yes {
				fmt.Printf("%s (killing)\n", key)
			}
		}
	case *listAll:
		for key, _ := range socialsMap {
			if state := configMap.data[key]; state {
				fmt.Printf("%s (killing)\n", key)
			} else {
				fmt.Println(key)
			}
		}
	case *add != "":
		if urls, ok := socialsMap[*add]; ok {
			configMap.add(*add)
			_ = configMap.save(configFile)
			writeEtcHosts(ETC_HOSTS, *add, urls)
			fmt.Printf("%s ghost is released\n", *add)
		} else {
			fmt.Println("no such ghost")
		}
	case *delete != "":
		if urls, ok := socialsMap[*delete]; ok {
			configMap.delete(*delete)
			_ = configMap.save(configFile)
			writeEtcHosts(ETC_HOSTS, *add, urls)
			fmt.Printf("%s ghost was killed\n", *add)
		} else {
			fmt.Println("no such ghost")
		}
	}
}

func init() {

	err := parseSocialsToMap(hostFile, socialsMap)
	if err != nil {
		log.Fatal(err)
	}

	for key, _ := range socialsMap {
		configMap.data[key] = false
	}
	err = parseConfigToMap(configFile, configMap.data)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
}

func writeEtcHosts(filename string, key string, urls []string) {
	var buf bytes.Buffer
	// write default hosts
	def, err := os.ReadFile(defaultGhosts)
	if err != nil {
		fmt.Println("Error reading %v: %v", defaultGhosts, err)
		return
	}
	_, err = buf.Write(def)
	if err != nil {
		fmt.Println("Error writing %v: %v", defaultGhosts, err)
	}

	// write un-released hosts
	_, err = buf.WriteString("\n# [" + key + "]" + "\n")
	if err != nil {
		fmt.Println("Error writing strings to buffer:", err)
		return
	}
	for _, url := range urls {
		_, err = buf.WriteString("0.0.0.0" + " " + url + "\n")
	}
	if err != nil {
		fmt.Println("Error writing strings to buffer:", err)
		return
	}

	// open /etc/hosts
	outputFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok && pathError.Err == syscall.EACCES {
			fmt.Println("Permission denied: You do not have write access to this file or directory.")
		} else {
			fmt.Println("Error opening file:", err)
		}
		return
	}
	defer outputFile.Close()

	// Use io.Copy to copy buffer content to the file
	_, err = io.Copy(outputFile, &buf)
	if err != nil {
		fmt.Println("Error copying buffer to file:", err)
		return
	}
}
