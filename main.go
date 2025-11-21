package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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
	SOCIAL_HOSTS  = "testdata/source.txt"
	DEFAULT_HOSTS = "testdata/default.txt"
	CONFIG_FILE   = "testdata/config"
	ETC_HOSTS     = "testdata/etc/hosts"
)

var (
	list = flag.Bool("list", false, "list all SM domains")
	add  = flag.String("add", "", "add a domain to Hosts")
	del  = flag.String("del", "", "delete a domain from Hosts")
)

var (
	ConfigMap = GhostMap{
		data: make(map[string][]string),
	}
	SocialMap = GhostMap{
		data: make(map[string][]string),
	}
)

func init() {
	// Social Media Map
	err := populateSocialMap(SOCIAL_HOSTS, SocialMap)
	if err != nil {
		log.Fatalf("Error parsing SM hosts file: %q", err)
	}
	// Config Map
	err = populateConfigMap(CONFIG_FILE, ConfigMap, SocialMap)
	if err != nil {
		log.Fatalf("Error parsing config file: %q", err)
	}
	flag.Parse()
}

func main() {
	switch {
	case *list:
		fmt.Println(printListView(SocialMap, ConfigMap))
	case *add != "":
		if yes := SocialMap.IsExists(*add); !yes {
			log.Fatalf("source hosts: %v, no such Host", *add)
		}
		if yes := ConfigMap.IsExists(*add); yes {
			fmt.Fprintf(os.Stdout, "config: %v, already added\n", *add)
			return
		}
		ConfigMap.data[*add] = SocialMap.data[*add]
		err := SaveConfigAndHosts()
		if err != nil {
			log.Fatal(err)
		}
	case *del != "":
		if yes := SocialMap.IsExists(*del); !yes {
			log.Fatalf("source hosts: %v, no such Host", *del)
		}
		if yes := ConfigMap.IsExists(*del); !yes {
			log.Fatalf("config: %v, no such Host", *del)
		}
		delete(ConfigMap.data, *del)
		err := SaveConfigAndHosts()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SaveConfigAndHosts() error {
	// write to config
	// write DEFAULT to hosts
	// write/append SOCIAL to hosts

	// For writing the Default hosts
	defHostFile, err := os.Open(DEFAULT_HOSTS)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer defHostFile.Close()

	etcHostFile, err := os.Create(ETC_HOSTS)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer etcHostFile.Close()

	writer := bufio.NewWriter(etcHostFile)

	_, err = io.Copy(writer, defHostFile)
	if err != nil {
		return fmt.Errorf("Error copying existing data: %v", err)
	}

	_, _ = writer.WriteString("\n##### GHOSTS #####\n")

	// append the SM hosts
	_, err = writer.WriteString(ConfigMap.String())
	if err != nil {
		return fmt.Errorf("Error writing new data: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error flusing data: %v", err)
	}
	fmt.Fprintln(os.Stdout, "Data successfully written to /etc/hosts")

	result := ConfigMap.List()
	err = os.WriteFile(CONFIG_FILE, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("Error writing config: %v", err)
	}
	fmt.Fprintln(os.Stdout, "Config successfully written to config")
	return nil
}
