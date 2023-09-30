package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

var (
	verbose = flag.Bool("v", false, "verbose")
)

// errandexit
func ErrAndExit(code int, err error) {
	log.Fatal(err)
}

var (
	taskFile = flag.String("t", "", "task file path")
)

func main() {
	flag.Parse()
	// global config from yaml file
	readGlobalConfigFromYaml(*configFilePath)
	// override global config from cli parameter
	readGlobalConfigFromCLI()

	// load inventory file
	hostConfs := ReadHostConfigFromYaml(*inventoryFilePath)
	if verbose != nil && *verbose {
		for _, host := range hostConfs {
			fmt.Println(host)
		}
	}

	t, err := LoadTask(*taskFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(hostConfs))
	for _, host := range hostConfs {
		go func(host *HostConfig) {
			hostConn := NewHostConn(host)
			err = hostConn.Exec(t)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(host)
	}
	wg.Wait()
}
