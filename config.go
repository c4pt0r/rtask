package main

import (
	"errors"
	"flag"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SshKeyPath  string `yaml:"ssh_key_path"`
	SshUser     string `yaml:"ssh_user"`
	SshPassword string `yaml:"ssh_password"`
}

type HostConfig struct {
	Host    string            `yaml:"host"`
	Port    int               `yaml:"port"`
	Timeout int               `yaml:"timeout"` // seconds
	Context map[string]string `yaml:"context"`
	// defualt config is global config
	SshConfig Config `yaml:"ssh_config"`
}

func (h *HostConfig) String() string {
	r, _ := yaml.Marshal(h)
	return string(r)
}

var (
	globalConfig Config
)

// read global config from cli parameter
var (
	configFilePath = flag.String("c", "", "config file path")
	sshKeyPath     = flag.String("ssh-key-path", "", "ssh key path")
	sshUser        = flag.String("ssh-user", "root", "ssh user")
	sshPwd         = flag.String("ssh-password", "", "ssh password")
)

func readGlobalConfigFromCLI() {
	if globalConfig.SshKeyPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			ErrAndExit(1, err)
		}
		if *sshKeyPath == "" {
			*sshKeyPath = path.Join(home, ".ssh/id_rsa")
		}
		globalConfig.SshKeyPath = *sshKeyPath
	}
	if globalConfig.SshUser == "" {
		globalConfig.SshUser = *sshUser
	}
	if globalConfig.SshPassword == "" {
		globalConfig.SshPassword = *sshPwd
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

// read inventory from yaml file
var (
	inventoryFilePath = flag.String("i", "", "inventory file path")
)

func ReadHostConfigFromYaml(inventoryFilePath string) []*HostConfig {
	if _, err := os.Stat(inventoryFilePath); os.IsNotExist(err) {
		ErrAndExit(1, errors.New("inventory file not found"))
	}
	// read config file
	configFile, err := os.ReadFile(inventoryFilePath)
	if err != nil {
		ErrAndExit(1, err)
	}
	// parse config file
	ret := []*HostConfig{}
	err = yaml.Unmarshal(configFile, &ret)
	if err != nil {
		ErrAndExit(1, err)
	}
	// make sure every host has a ssh config
	for _, hostConfig := range ret {
		if hostConfig.Port == 0 {
			hostConfig.Port = 22
		}
		if hostConfig.Timeout == 0 {
			hostConfig.Timeout = 10
		}
		if hostConfig.SshConfig.SshKeyPath == "" {
			hostConfig.SshConfig.SshKeyPath = globalConfig.SshKeyPath
		}
		if hostConfig.SshConfig.SshUser == "" {
			hostConfig.SshConfig.SshUser = globalConfig.SshUser
		}
		if hostConfig.SshConfig.SshPassword == "" {
			hostConfig.SshConfig.SshPassword = globalConfig.SshPassword
		}
	}
	return ret
}
