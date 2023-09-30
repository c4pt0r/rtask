package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Task struct {
	Name string `yaml:"name"`
	Tp   string `yaml:"type"`
	/*
		CopyTask struct {
			Local  string `yaml:"local_path"`
			Remote string `yaml:"remote_path"`
		}
	*/
	ExecTask struct {
		Cmd string `yaml:"cmd"`
	} `yaml:"exec"`

	Timeout int `yaml:"timeout"` // seconds
}

func (t *Task) IsExecTask() bool {
	return t.Tp == "exec"
}

func render(v map[string]string, tplStr string) (string, error) {
	tpl, err := template.New("cmd").Parse(tplStr)
	if err != nil {
		return "", err
	}
	stringBuffer := &strings.Builder{}
	err = tpl.Execute(stringBuffer, v)
	if err != nil {
		return "", err
	}
	return stringBuffer.String(), nil
}

func (t *Task) GetCommand(ctx context.Context) (string, error) {
	if !t.IsExecTask() {
		return "", errors.New("not exec task")
	}
	ret := t.ExecTask.Cmd
	var err error
	if v, ok := ctx.Value("context").(map[string]string); ok {
		ret, err = render(v, ret)
		if err != nil {
			return "", err
		}
	}
	return ret, nil
}

func LoadTask(filePath string) (*Task, error) {
	// read yaml file
	task := &Task{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ErrAndExit(1, errors.New("task file not found"))
	}
	// read config file
	taskFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(taskFile, &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
