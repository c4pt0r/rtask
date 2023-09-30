package main

import (
	"context"
	"fmt"
	"time"

	"github.com/appleboy/easyssh-proxy"
)

type HostConn struct {
	hostConfig *HostConfig

	ssh *easyssh.MakeConfig

	stdoutChan <-chan string
	stderrChan <-chan string
	doneChan   <-chan bool
	errChan    <-chan error
}

func NewHostConn(hostConfig *HostConfig) *HostConn {
	c := &HostConn{hostConfig: hostConfig}
	c.ssh = &easyssh.MakeConfig{
		Server:  c.hostConfig.Host,
		User:    c.hostConfig.SshConfig.SshUser,
		Port:    fmt.Sprintf("%d", c.hostConfig.Port),
		Timeout: time.Duration(c.hostConfig.Timeout) * time.Second,
	}
	if c.hostConfig.SshConfig.SshPassword != "" {
		c.ssh.Password = c.hostConfig.SshConfig.SshPassword
	} else {
		c.ssh.KeyPath = c.hostConfig.SshConfig.SshKeyPath
	}
	return c
}

func now() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func (c *HostConn) Exec(t *Task) error {
	var err error
	// build context

	ctx := context.WithValue(context.Background(), "context", c.hostConfig.Context)
	cmd, err := t.GetCommand(ctx)
	if err != nil {
		return err
	}
	c.stdoutChan, c.stderrChan, c.doneChan, c.errChan, err = c.ssh.Stream(cmd, time.Duration(c.hostConfig.Timeout)*time.Second)
	// Handle errors
	if err != nil {
		return err
	} else {
		// read from the output channel until the done signal is passed
		isTimeout := true
	loop:
		for {
			select {
			case isTimeout = <-c.doneChan:
				fmt.Printf("[%s | FINISHED | %s]\n", c.hostConfig.Host, now())
				break loop
			case outline := <-c.stdoutChan:
				if outline != "" {
					fmt.Printf("[%s | STDOUT | %s] %s\n", c.hostConfig.Host, now(), outline)
				}
			case errline := <-c.stderrChan:
				if errline != "" {
					fmt.Printf("[%s | STDERR | %s] %s\n", c.hostConfig.Host, now(), errline)
				}
			case err = <-c.errChan:
			}
		}
		// get exit code or command error.
		if err != nil {
			fmt.Printf("[%s | ERR | %s] %s\n", c.hostConfig.Host, now(), err.Error())
		}

		// command time out
		if !isTimeout {
			fmt.Printf("[%s | TIMEOUT | %s]\n", c.hostConfig.Host, now())
		}
	}
	return nil
}
