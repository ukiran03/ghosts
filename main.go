package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Reset = "\033[0m"
)

var (
	SOCIAL_HOSTS  = "/home/ukiran/.config/ghosts/socials.hosts"
	DEFAULT_HOSTS = "/home/ukiran/.config/ghosts/default.hosts"
	CONFIG_FILE   = "/home/ukiran/.config/ghosts/config"
	ETC_HOSTS     = "/etc/hosts" // root
)

const Help = `
ghosts [option] [flag] [args]
    add:    add a site to the Hosts
            '--all' flag adds all the sites
    del:    delete a site from the Hosts
            '--all' flag deletes all sites
    list:   lists the sites that are added and deleted
    help:   help that you see`

var (
	SocialMap = GhostMap{data: make(map[string][]string)}
	ConfigMap = GhostMap{data: make(map[string][]string)}
)

var (
	addCmd = flag.NewFlagSet("add", flag.ExitOnError)
	addAll = addCmd.Bool("all", false, "add all sites to Hosts")
	delCmd = flag.NewFlagSet("del", flag.ExitOnError)
	delAll = delCmd.Bool("all", false, "del all sites from Hosts")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: \t%v\n", Help)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addAll {
			maps.Copy(ConfigMap.data, SocialMap.data)
		} else {
			args := addCmd.Args()
			if len(args) != 0 {
				for _, arg := range args {
					if yes := SocialMap.IsExists(arg); !yes {
						log.Fatalf("source hosts: %v, no such Host\n", arg)
					}
					if yes := ConfigMap.IsExists(arg); yes {
						fmt.Fprintf(os.Stdout, "config: %v, already added\n", arg)
					}
					ConfigMap.data[arg] = SocialMap.data[arg]
				}
			} else {
				fmt.Fprintf(os.Stderr, "Given No argument\n")
				os.Exit(1)
			}
		}
		err := SaveConfigAndHosts()
		if err != nil {
			log.Fatal(err)
		}

	case "del":
		delCmd.Parse(os.Args[2:])
		if *delAll {
			ConfigMap.data = make(map[string][]string) // Reset
		} else {
			args := delCmd.Args()
			if len(args) != 0 {
				for _, arg := range args {
					if yes := SocialMap.IsExists(arg); !yes {
						log.Fatalf("source hosts: %v, no such Host", arg)
					}
					if yes := ConfigMap.IsExists(arg); !yes {
						log.Fatalf("config: %v, no such Host", arg)
					}
					delete(ConfigMap.data, arg)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Given No argument\n")
				os.Exit(1)
			}
		}
		err := SaveConfigAndHosts()
		if err != nil {
			log.Fatal(err)
		}
	case "list":
		fmt.Println(printListView(SocialMap, ConfigMap))
	case "help":
		fmt.Printf("Usage: \t%v\n", Help)
	default:
		fmt.Fprintf(os.Stderr, "Usage: \t%v\n", Help)
	}
}

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

func SaveConfigAndHosts() error {
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

	// Append the SM hosts
	_, err = writer.WriteString(ConfigMap.String())
	if err != nil {
		return fmt.Errorf("Error writing new data: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error flusing data: %v", err)
	}

	// Write to Config
	configFile, err := os.Create(CONFIG_FILE)
	if err != nil {
		return fmt.Errorf("Error creating config file: %v", err)
	}
	defer configFile.Close()
	configData := ConfigMap.List()
	_, err = configFile.WriteString(configData)
	if err != nil {
		return fmt.Errorf("Error writing config file: %v", err)
	}
	return nil
}
