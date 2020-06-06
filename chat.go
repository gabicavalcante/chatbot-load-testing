package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"bytes"
)
 

func getRule(body string) (rule ChatRule, ok bool) {
	for _, r := range config.Rules { 
		if strings.Contains(body, r.Receive.Content) {
			rule, ok = r, true
			return
		}
	}
	return
} 

func sendResponse(rule ChatRule, reqBody_ string) { 
	if config.ResponsePauseMax > 0 {
		ms := rand.Intn(config.ResponsePauseMax - config.ResponsePauseMin)
		ms += config.ResponsePauseMin
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	body := rule.Send.Content
	data := &Data{Talk_id: config.UserId, Clinic_id: 63, Messages: []string{body}}
    reqBody, _ := json.Marshal(data) 

	res, err := http.Post(
		config.ChatBotURL,
		"application/json",
        bytes.NewBuffer(reqBody))
        
	if err != nil {
		fmt.Println(err)
    } else {
		if res.StatusCode == 200 {
			fmt.Printf("sucess: %s", rule.Send.Name)
		} else {
			fmt.Printf("fail: %s", rule.Send.Name)
		}
	}
}


func readRequest(r *http.Request) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	in := string(bodyBytes)

	var data map[string]interface{}
	json.Unmarshal([]byte(in), &data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%v] %v", r.Method, r.URL) 
	fmt.Printf("> talk_id: %f, thread_count: %f, cpu: %f, ram: %f\n", data["talk_id"], data["thread_count"], data["cpu"], data["ram"]) 


    if config.RequestTimeMax > 0 {
		ms := rand.Intn(config.RequestTimeMax - config.RequestTimeMin)
		ms += config.RequestTimeMin
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	rule, ok := getRule(fmt.Sprint(data["text"]))
	if !ok {
		fmt.Printf("rule: %s\n",  "OTHERS")
		return nil
	}
	
	fmt.Printf("rule: %s\n",  rule.Send.Name) 

	go sendResponse(rule, fmt.Sprint(data["text"]))

	return nil
}