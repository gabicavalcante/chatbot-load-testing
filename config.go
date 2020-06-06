package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ChatBotURL       string
	RequestTimeMin   int // min time for processing request, ms
	RequestTimeMax   int
	ResponsePauseMin int // min time for thinking on the answer, ms
	ResponsePauseMax int
	UserId			 int
	Rules          []ChatRule
}

type ChatRule struct {
	Send  		ChatRequest
	Receive 	ChatResponse
}

type ChatRequest struct {
	Name        string 
	Content		string // body should contain the segment
}

type ChatResponse struct { 
	Content		string
}

func (conf *Config) Prepare() {
	for i, rule := range conf.Rules {
		if rule.Send.Name == "" {
			conf.Rules[i].Send.Name = fmt.Sprintf("[%v])",
				rule.Send.Content)
		}
	}
}

func readConfig(config *Config) {
	blob, err := ioutil.ReadFile(*configFile)
	if err != nil || len(blob) == 0 {
		fmt.Println("error:", err)
		fmt.Println("Can't read config file. You should use \"-conf\" flag.")
		os.Exit(0)
	}

	err = json.Unmarshal(blob, config)
	if err != nil {
		fmt.Println("error:", err)
		fmt.Printf("Can't decode config file [%v]\n", *configFile)
		os.Exit(0)
	}

	config.Prepare()
}
