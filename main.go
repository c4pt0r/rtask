package main

import (
	"flag"
	"fmt"
)

// errandexit
func ErrAndExit(code int, err error) {
	panic(err)
}

var (
	taskFile = flag.String("task", "", "task file path")
)

func main() {
	flag.Parse()
	// global config from yaml file
	readGlobalConfigFromYaml(*configFilePath)
	// override global config from cli parameter
	readGlobalConfigFromCLI()

	// load inventory file
	hostConfs := ReadHostConfigFromYaml(*inventoryFilePath)
	for _, host := range hostConfs {
		fmt.Println(host)
	}

	t, err := LoadTask(*taskFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, host := range hostConfs {
		hostConn := NewHostConn(host)
		err = hostConn.Exec(t)
		if err != nil {
			fmt.Println(err)
		}
	}
}
