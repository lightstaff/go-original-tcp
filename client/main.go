package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lightstaff/go-original-tcp/protocol"
)

var (
	message = flag.String("message", "", "send message")
	loop    = flag.Int("loop", 1, "loop count")
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	for i := 0; i < *loop; i++ {
		wg.Add(1)
		go func(i int, s int64) {
			defer func() {
				wg.Done()
				log.Printf("[INFO] Tick: %d\n", time.Now().UnixNano()-s)
			}()

			conn, err := net.Dial("tcp", "localhost:18888")
			if err != nil {
				log.Printf("[ERROR] %s\n", err.Error())
				return
			}
			defer conn.Close()

			m := &protocol.Message{
				Headers: map[string]string{"Count": fmt.Sprintf("%d", i)},
				Body:    *message,
			}

			if err := m.Write(conn); err != nil {
				log.Printf("[ERROR] %s\n", err.Error())
				return
			}

			log.Printf("[INFO] Message %v\n", m)

			r, err := protocol.ReadReply(conn)
			if err != nil {
				log.Printf("[ERROR] %s\n", err.Error())
				return
			}

			log.Printf("[INFO] Reply %v\n", r)
		}(i, time.Now().UnixNano())
	}

	wg.Wait()
}
