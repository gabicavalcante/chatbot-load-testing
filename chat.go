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
		if strings.Contains(body, r.Response.Content) {
			rule, ok = r, true
			return
		}
	}
	return
} 

func sendResponse(rule ChatRule) { 
	statChan <- Stat{Name: rule.Request.Name, ReqCount: 1, Request: time.Now()}
	fmt.Printf("sedings... %s\n", rule.Request.Name)

	if config.ResponsePauseMax > 0 {
		ms := rand.Intn(config.ResponsePauseMax - config.ResponsePauseMin)
		ms += config.ResponsePauseMin
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	body := rule.Request.Content

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
			fmt.Printf("sucess: %s\n", rule.Request.Name)
		} else {
			fmt.Printf("fail: %s\n", rule.Request.Name)
			statChan <- Stat{Name: rule.Request.Name, RespCount: 1, Response: time.Now()}
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
		statChan <- Stat{Name: otherRequests, ReqCount: 1, Request: time.Now()}
		return nil
	}

	statChan <- Stat{Name: rule.Request.Name, RespCount: 1, TalkId: data["talk_id"].(float64), CPU: data["cpu"].(float64), ThreadCount: data["thread_count"].(float64), RAM: data["ram"].(float64), Response: time.Now()}
	
	fmt.Printf("rule: %s. content: %s\n",  rule.Request.Name, rule.Request.Content) 

	if rule.Request.Content != "" {
		go sendResponse(rule)
	} else {
		fmt.Printf("-- end\n")
	}
	return nil
}