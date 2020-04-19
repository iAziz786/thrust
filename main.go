package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
)

var MAX int

type MyEvent struct {
	Name string `json:"name"`
}

// func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
// 	defer cancel()
// 	sem := make(chan int, MAX)
// 	i := 0
// 	sigs := make(chan os.Signal)
// 	done := make(chan struct{})
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// 	go func() {
// 		<-sigs
// 		fmt.Println("closing signal received")
// 		close(done)
// 	}()
// Outerloop:
// 	for {
// 		select {
// 		case <-done:
// 			fmt.Println("closed done")
// 			break Outerloop
// 		case <-ctx.Done():
// 			break Outerloop
// 		case <-time.After(30 * time.Second):
// 			break Outerloop
// 		default:
// 			i++
// 			sem <- 1 // will block if there is MAX ints in sem
// 			go func(count int) {
// 				fmt.Println("request", count)
// 				resp, err := http.Get("http://45.76.160.109/")
// 				if err != nil {
// 					<-sem
// 					return
// 				}
// 				fmt.Println("receiving...")
// 				resp.Body.Close()
// 				<-sem // removes an int from sem, allowing another to proceed
// 			}(i)
// 		}
// 	}
// 	return fmt.Sprintf("Done"), nil
// }

type Value struct {
	value   int
	latency time.Duration
	lock    sync.Mutex
}

func main() {
	// lambda.Start(HandleRequest)
	url := os.Args[3]
	MAX, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	fmt.Println(MAX)
	sem := make(chan struct{}, MAX)
	value := Value{}
	sigs := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("closing signal received")
		close(done)
	}()
	elapse := os.Args[1]
	go func() {
		second, err := strconv.Atoi(elapse)
		if err != nil {
			second = 10
		}
		<-time.After(time.Duration(second) * time.Second)
		close(done)
	}()
	now := time.Now()
Outerloop:
	for {
		select {
		case <-done:
			fmt.Println("closed done")
			break Outerloop
		default:
			sem <- struct{}{} // will block if there is MAX ints in sem
			go func() {
				select {
				case <-done:
					return
				default:
					startTime := time.Now()
					resp, err := http.Get(url)
					if err != nil {
						<-sem
						return
					}
					value.lock.Lock()
					value.value++
					fmt.Printf("#%d Received\n", value.value)
					value.latency += time.Since(startTime)
					value.lock.Unlock()
					resp.Body.Close()
					<-sem // removes an int from sem, allowing another to proceed
				}
			}()
		}
	}
	fmt.Printf("%d request completed in %s\n", value.value, time.Since(now))
	fmt.Printf("Average latency: %s\n", (value.latency / time.Duration(value.value)))
}

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://4023445056dd46aea2ce8e6bb99ca997@sentry.io/5178758",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
