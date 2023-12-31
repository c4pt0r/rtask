package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

var (
	printSampleTask        = flag.Bool("print-sample-task", false, "print sample task file")
	sampleTaskFile  string = `name: "simple echo and sleep" 
type: "exec"
exec:
  cmd: |
    echo "read context from inventory.yaml for different nodes"
    echo Hello on {{ .HOSTNAME }} > /tmp/hello_from_rsync
    echo "start sleeping..."
    sleep {{ .SLEEP }}`

	printSampleInventory        = flag.Bool("print-sample-inventory", false, "print sample inventory file")
	sampleInventoryFile  string = `- host: "node1"
  port: 22
  timeout: 10
  context:
    HOSTNAME: "node1"
    SLEEP: "3"

- host: "node2"
  port: 22
  timeout: 10 
  context:
    HOSTNAME: "node2"
    SLEEP: "5"
`
	printSampleConfig        = flag.Bool("print-sample-config", false, "print sample config file")
	sampleConfigFile  string = `ssh_key_path: "/home/root/.ssh/id_rsa"
ssh_user: "root"`
)

func doPrintSampleTaskFile() {
	fmt.Println(sampleTaskFile)
}

func doPrintSampleInventoryFile() {
	fmt.Println(sampleInventoryFile)
}

func doPrintSampleConfigFile() {
	fmt.Println(sampleConfigFile)
}

func printSample() {
	if *printSampleTask {
		doPrintSampleTaskFile()
		os.Exit(0)
	}
	if *printSampleInventory {
		doPrintSampleInventoryFile()
		os.Exit(0)
	}
	if *printSampleConfig {
		doPrintSampleConfigFile()
		os.Exit(0)
	}
	return
}

func main() {
	flag.Parse()
	printSample()

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

	// TODO: use a pool to limit the number of concurrent connections
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
