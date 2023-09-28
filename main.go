package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	SshKeyPath string `yaml:"ssh_key_path"`
	User       string `yaml:"user"`
}

type TaskConfig struct {
	TaskName string `yaml:"task_name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	KeyPath  string `yaml:"ssh_key_path"`
	AsUser   string `yaml:"as_user"`
	Timeout  int    `yaml:"timeout"`

	Command string            `yaml:"command"`
	Context map[string]string `yaml:"context"`
}

var (
	globalConfig GlobalConfig
)

// errandexit
func ErrAndExit(code int, err error) {
	fmt.Println(err)
	os.Exit(code)
}

// read global config from cli parameter
var (
	configFilePath = flag.String("config", "", "config file path")
)

func readGlobalConfigFromCLI() {
	if globalConfig.SshKeyPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			ErrAndExit(1, err)
		}
		sshPath := path.Join(home, ".ssh")
		globalConfig.SshKeyPath = *flag.String("ssh-key-path", sshPath, "ssh key path")
	}
	if globalConfig.User == "" {
		// get current username
		currentUser := os.Getenv("USER")
		if currentUser == "" {
			currentUser = "root"
		}
		globalConfig.User = *flag.String("user", currentUser, "ssh user")
	}
}

func readGlobalConfigFromYaml(configFilePath string) {
	// if there's no config file, just return
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return
	}
	// read config file
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		ErrAndExit(1, err)
	}
	// parse config file
	err = yaml.Unmarshal(configFile, &globalConfig)
	if err != nil {
		ErrAndExit(1, err)
	}
}

func init() {
	flag.Parse()
	// override global config from yaml file
	readGlobalConfigFromYaml(*configFilePath)
	// read global config from cli parameter
	readGlobalConfigFromCLI()
}

func main() {
	// Create MakeConfig instance with remote username, server address and path to private key.
	/*
		ssh := &easyssh.MakeConfig{
			Server:  "localhost",
			User:    "drone-scp",
			KeyPath: "./tests/.ssh/id_rsa",
			Port:    "22",
			Timeout: 60 * time.Second,
		}

		// Call Run method with command you want to run on remote server.
		stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream("for i in {1..5}; do echo ${i}; sleep 1; done; exit 2;", 60*time.Second)
		// Handle errors
		if err != nil {
			panic("Can't run remote command: " + err.Error())
		} else {
			// read from the output channel until the done signal is passed
			isTimeout := true
		loop:
			for {
				select {
				case isTimeout = <-doneChan:
					break loop
				case outline := <-stdoutChan:
					fmt.Println("out:", outline)
				case errline := <-stderrChan:
					fmt.Println("err:", errline)
				case err = <-errChan:
				}
			}

			// get exit code or command error.
			if err != nil {
				fmt.Println("err: " + err.Error())
			}

			// command time out
			if !isTimeout {
				fmt.Println("Error: command timeout")
			}
		}
	*/

	// just dump the Config for debugging
	fmt.Println(globalConfig)
}
