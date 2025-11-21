package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Reset = "\033[0m"
)

var (
	SOURCE_FILE = "testdata/source.txt"
	CONFIG_FILE = "testdata/config"
	HOST_FILE   = "testdata/etc/hosts"
)

var (
	list = flag.Bool("list", false, "list all SM domains")
	add  = flag.String("add", "", "add a domain to Hosts")
	del  = flag.String("del", "", "delete a domain from Hosts")
)

var (
	Cmap = GhostMap{
		data: make(map[string][]string),
	}
	SMmap = GhostMap{
		data: make(map[string][]string),
	}
)

func init() {
	// Social Media Map
	err := populateSocialMap(SOURCE_FILE, SMmap)
	if err != nil {
		log.Fatalf("Error parsing SM hosts file: %q", err)
	}
	// Config Map
	err = populateConfigMap(CONFIG_FILE, Cmap, SMmap)
	if err != nil {
		log.Fatalf("Error parsing config file: %q", err)
	}
	flag.Parse()
}

func main() {
	switch {
	case *list:
		fmt.Println(printListView(SMmap, Cmap))
	case *add != "":
		if yes := SMmap.IsExists(*add); !yes {
			log.Fatalf("source hosts: %v, no such Host", *add)
		}
		if yes := Cmap.IsExists(*add); yes {
			fmt.Fprintf(os.Stdout, "config: %v, already added\n", *add)
			return
		}
		Cmap.data[*add] = SMmap.data[*add]
		err := saveConfig()
		if err != nil {
			log.Fatal(err)
		}
	case *del != "":
		if yes := SMmap.IsExists(*del); !yes {
			log.Fatalf("source hosts: %v, no such Host", *del)
		}
		if yes := Cmap.IsExists(*del); !yes {
			log.Fatalf("config: %v, no such Host", *del)
		}
		delete(Cmap.data, *del)
		err := saveConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveConfig() error {
	result := Cmap.List()
	err := Cmap.SaveToFile(HOST_FILE)
	if err != nil {
		return err
	}
	return os.WriteFile(CONFIG_FILE, []byte(result), 0644)
}
