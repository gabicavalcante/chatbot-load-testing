package main

import (
	"encoding/json"
    "flag"
    "fmt"
    "net/http"
    "time"
    "io/ioutil"
    "log"
    "sync"
    "bytes"
    "os"
)

var portNumber *string
var configFile *string
var config Config

type Data struct{
    Talk_id     int         `json:"talk_id"`
    Clinic_id   int         `json:"clinic_id"`
    Messages    []string    `json:"messages"`
}


func botsHandler(w http.ResponseWriter, r *http.Request) {
    err := readRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}


func statHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("[%v] %v\n", r.Method, r.URL)
    printStat(w) 
    printStat(os.Stdout)
}


func MakeRequest(url string, talk_id int, ch chan<-string) {
    start := time.Now()

    data := &Data{Talk_id: talk_id, Clinic_id: 63, Messages: []string{"marcar consulta"}}
    reqBody, _ := json.Marshal(data) 

	resp, _ := http.Post(url, "application/json;", bytes.NewBuffer(reqBody))

    secs := time.Since(start).Seconds()

    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    ch <- fmt.Sprintf("%.2f elapsed with response length: %d %d %s\n", secs, len(bodyBytes), resp.StatusCode, url)
    ch <- fmt.Sprintf("> %s \n", string(bodyBytes))
}

func main() {
    portNumber = flag.String("port", "8877", "Port number to use for connection")
    configFile = flag.String("conf", "config.json", "Path to config file")
    flag.Parse()
    
    readConfig(&config)
    initStat(config)

    start := time.Now()
    ch := make(chan string)

    // create a WaitGroup
    wg := new(sync.WaitGroup)
    // add two goroutines to `wg` WaitGroup
    wg.Add(2)

    http.HandleFunc("/stat", statHandler)
    http.HandleFunc("/message", botsHandler)

    // create a default route handler
    http.HandleFunc( "/", func( res http.ResponseWriter, req *http.Request ) {
        fmt.Fprint( res, "Hello: " + req.Host )
    } )

    go func() {
        fmt.Println("server up!")
        log.Fatal(http.ListenAndServe(":"+*portNumber, nil))
        wg.Done() // one goroutine finished
    }()

    // goroutine to launch a server on port 9000
    go func() {
        log.Fatal( http.ListenAndServe( ":9000", nil ) )
        wg.Done() // one goroutine finished
    }()
     
    for i := 0; i < 1; i++ {
        go MakeRequest(config.ChatBotURL, config.UserId, ch)
    }

    for i := 0; i < 1; i++ {
        fmt.Println(<-ch)
    }
    fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

    // wait until WaitGroup is done
    wg.Wait()
}
