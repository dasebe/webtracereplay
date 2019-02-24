package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"bufio"
	"os"
	"time"
	"sync"
	"net/http"
	"strings"   
	"sync/atomic"
)
var trace *string = flag.String("clientTrace", "w500m.tr", "client requesting trace")
var numThread *int = flag.Int("numThread", 30, "number of thread that client use")
var ip *string = flag.String("ip", "127.0.0.1:6081", "ip to connect")
var throughputLog *string = flag.String("throughputLog", "throughput.log", "throughput log")
var histogramLog *string = flag.String("histogramLog", "histogram.log", "histogram log")

type Empty struct{}
type Semaphore chan Empty
func NewSemaphore(n int) Semaphore {
    return make(Semaphore, n)
}
func (s Semaphore) Acquire() {
    s <- Empty{}
}
func (s Semaphore) Release() {
    <- s
}

var httpSema Semaphore

// error check alias
func Check(err error) {
    if err != nil {
        fmt.Println("error")
        panic(err)
    }
}

var schedCount uint64 = 0
var sentCount uint64 = 0
var printCount int = 0
var size uint64 = 0
var received uint64 = 0
var sent uint64 = 0
var setZero uint64 = 0

// http semaphore init
var clientPool *http.Client
func InitHttp(numThread int) {
	httpSema = NewSemaphore(numThread)
    fmt.Println("init begin")
    dfTransp := http.DefaultTransport
    dfTranspPtr, ok := dfTransp.(*http.Transport)
    if !ok {
        panic(fmt.Sprintf("failed to initialize http transport"))
    }
    dfTranspStruct := *dfTranspPtr
    dfTranspStruct.MaxIdleConns = numThread
    dfTranspStruct.MaxIdleConnsPerHost = numThread
    clientPool = &http.Client{Transport: &dfTranspStruct}
    fmt.Println("init end")
}

// send json request to app server
func SendRequest(s string) ([]byte, int) {
	// execute query
	httpSema.Acquire()
	id := strings.Split(s, " ")
	url := "http://" + *ip + "/" + id[1];
	resp, err := clientPool.Get(url)
	httpSema.Release()
	atomic.AddUint64(&sentCount, 1)
	if err != nil{
		fmt.Println("error issueing post",s)
		panic(err)
	}
	val, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	atomic.AddUint64(&received, 1)
	if err != nil{
		fmt.Println("error reading resp",s)
		filename := "./error.dat"
		f, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		defer f.Close()
		f.WriteString("read body error.\n")
	}
	return val, 1
}

// read json trace file and encode requests
func LoadGenerator(trace string) {
	var wg sync.WaitGroup
	f, err := os.Open(trace)
	Check(err)
	scanner := bufio.NewScanner(f)
	// Main loop
	printCounter := 0
	for scanner.Scan() {
		for schedCount - sentCount > 1000 {
			printCount++
			if printCount>100 {
	    		fmt.Println("Closed loop (100x)",time.Now())
	    		printCount = 0
			}
			time.Sleep(time.Millisecond*10)
		}
		// issue request
		wg.Add(1)
		atomic.AddUint64(&schedCount, 1)
		atomic.AddUint64(&sent, 1)
		go func(body string) {
			if len(body) > 0 {
				SendRequest(body)
			}
			wg.Done()
		} (scanner.Text())

		if printCounter> 5000 {
			printCounter = 0
			fmt.Println("Status",schedCount, sentCount, "sema", len(httpSema), cap(httpSema))
		}
		time.Sleep(time.Microsecond*50)
		printCounter++
	}
	// Wait for completion
	fmt.Println("done sending", schedCount, sentCount, "sema", len(httpSema), cap(httpSema))
	wg.Wait()
	fmt.Println("done waiting", schedCount, sentCount, "sema", len(httpSema), cap(httpSema))
}

func throughput() {
	path := "/home/wenqim/result/clientThroughput"
	_ = os.Remove(path)
	_, err := os.Create(path)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	Check(err)
	defer f.Close()
	for{
		select{
			case <-time.After(1 * time.Second):
				str1 := fmt.Sprintf("sent: %d/s\n", sent)
				str2 := fmt.Sprintf("received: %d/s\n", received)
				_, err = f.WriteString(str1)
				_, err = f.WriteString(str2)
				Check(err)
				received = atomic.LoadUint64(&setZero)
				sent = atomic.LoadUint64(&setZero)
		}
	}
}
func main() {
    flag.Parse()
    fmt.Println("begin")
	InitHttp(*numThread)
	go throughput()
	LoadGenerator(*trace)
    fmt.Println("done")
}
