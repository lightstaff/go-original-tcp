package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/lightstaff/go-original-tcp/protocol"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	listener, err := net.Listen("tcp", "localhost:18888")
	if err != nil {
		log.Fatalf("[ERROR] %s\n", err.Error())
	}
	defer listener.Close()

	go func(l net.Listener) {
	LISTENER_FOR:
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("[ERROR] %s\n", err.Error())
			}

			log.Printf("[INFO] Accepted %v\n", conn.RemoteAddr())

			go func(c net.Conn, s int64) {
				defer func() {
					conn.Close()

					log.Printf("[INFO] Tick: %d\n", time.Now().UnixNano()-s)
				}()

				m, err := protocol.ReadMessage(bufio.NewReader(c))
				if err == io.EOF {
					log.Println("[INFO] Connection closed")
					return
				}
				if err != nil {
					r := &protocol.Reply{
						Status: protocol.ERR,
					}

					if err := r.Write(c); err != nil {
						log.Printf("[ERROR] %s\n", err.Error())
					}

					log.Printf("[INFO] Reply %v\n", r)
					log.Printf("[ERROR] %s\n", err.Error())
					return
				}

				log.Printf("[INFO] Message %v\n", m)

				r := &protocol.Reply{
					Status: protocol.OK,
				}

				if err := r.Write(c); err != nil {
					log.Printf("[ERROR] %s\n", err.Error())
				}

				log.Printf("[INFO] Reply %v\n", r)
			}(conn, time.Now().UnixNano())

			select {
			case <-ctx.Done():
				break LISTENER_FOR
			default:
				continue
			}
		}
	}(listener)

	log.Println("[INFO] Start")

	<-signals
	cancel()

	log.Println("[INFO] Stop")
}
