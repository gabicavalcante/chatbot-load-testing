package main

import (
	"fmt"
	"io"
	"time"

	"github.com/olekukonko/tablewriter"
)

const otherRequests = "OTHERS"

type Stat struct {
	Name        string 
	ReqCount    int
	RespCount   int
	TalkId		float64
	CPU			float64
	ThreadCount	float64
	RAM			float64
	Request 	time.Time 
	Response 	time.Time 
}	

var statOrder []string
var statData map[string]*Stat
var statChan chan Stat

func (s Stat) getTitles() []string {
	return []string{
		"Name",  "Req Count", "Resp Count", "Talk Id", 
		"CPU", "Thread Count", "RAM", "Request", "Response",
	}
}

func (s Stat) getStrings() []string {
	return []string{
		s.Name,
		fmt.Sprint(s.ReqCount),
		fmt.Sprint(s.RespCount),
		fmt.Sprint(s.TalkId),
		fmt.Sprint(s.CPU),
		fmt.Sprint(s.ThreadCount),
		fmt.Sprint(s.RAM),
		fmt.Sprint(s.Request.Format("15:04:05.000")),
		fmt.Sprint(s.Response.Format("15:04:05.000")),
	}
}

func printStat(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(Stat{}.getTitles())

	for _, name := range statOrder {
		table.Append(statData[name].getStrings())
	}

	table.Render()
}

func statWriter() {
	for {
		s := <-statChan
		statData[s.Name].ReqCount += s.ReqCount
		statData[s.Name].RespCount += s.RespCount
		statData[s.Name].TalkId = s.TalkId 
		statData[s.Name].CPU = s.CPU 
		statData[s.Name].ThreadCount = s.ThreadCount 
		statData[s.Name].RAM = s.RAM
		statData[s.Name].Request = s.Request 
		statData[s.Name].Response = s.Response 
	}
}

func resetStat(config Config) {
	for _, rule := range config.Rules {
		statData[rule.Request.Name] = &Stat{Name: rule.Request.Name}
	}
	statData[otherRequests] = &Stat{Name: otherRequests}
}

func initStat(config Config) {
	statOrder = make([]string, len(config.Rules)+1)
	statData = make(map[string]*Stat)
	statChan = make(chan Stat, 1000)

	fmt.Printf(">> %d\n", len(config.Rules)+1)

	for i, rule := range config.Rules {
		statOrder[i] = rule.Request.Name
		statData[rule.Request.Name] = &Stat{Name: rule.Request.Name}
	}
	statOrder[len(config.Rules)] = otherRequests
	statData[otherRequests] = &Stat{Name: otherRequests}

	go statWriter()
}
