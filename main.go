// The Telnet client
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	host    = kingpin.Flag("host", "Host").Required().String()
	port    = kingpin.Flag("port", "Port").Required().String()
	timeout = kingpin.Flag("timeout", "Timeout on connection").Default(
		strconv.FormatInt(60, 10)).Int()
	requestTimeout = kingpin.Flag("request_timeout", "Timeout for each request").Default(
		strconv.FormatInt(1, 10)).Int()
)

func applyCommand(conn net.Conn, command string) {
	_, err := conn.Write([]byte(command))
	if err != nil {
		log.Fatalln(err)
	}

	err = conn.SetReadDeadline(
		time.Now().Add(time.Duration(*requestTimeout) * time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	b := make([]byte, 4096)
	for {
		n, err := conn.Read(b)
		if err != nil || n == 0 {
			break
		}
		fmt.Println(string(b[:n]))
	}
}

func makeReadChannel() <-chan string {
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text() + "\n"
		}
	}()
	return ch
}

func runUntilComplete() {
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	cmd := makeReadChannel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigs:
			return
		case command := <-cmd:
			applyCommand(conn, command)
		}
	}
}

func main() {
	kingpin.Parse()
	runUntilComplete()
}
