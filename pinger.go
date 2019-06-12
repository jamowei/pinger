package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("pinger", "Pings a running pinger server or starts a server to look for open tcp ports")

	var host = parser.String("s", "server", &argparse.Options{Required: false, Help: "host to connect with (server mode)"})
	var portRange = parser.String("r", "range", &argparse.Options{
		Required: false,
		Validate: func(args []string) error {
			if len(args) == 1 {
				_, err := extractRangeFromParam(args[0])
				return err
			}
			return errors.New("Error: range parameter has not enough arguments")
		},
		Help: "range of port numbers, e.g. 8080-8090",
	})
	var ports = parser.List("p", "port", &argparse.Options{Required: false, Help: "ports to listen on (server mode) or to connect with (client mode)"})

	logger := log.New(os.Stdout, "pinger - ", 0)

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if len(*ports) == 0 && len(*portRange) > 0 {
		ports, err = extractRangeFromParam(*portRange)
		if err != nil {
			logger.Fatal(err)
			os.Exit(1)
		}
	}

	if len(*ports) == 0 {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if len(*host) > 0 { // client mode
		clientMode(*host, *ports, logger)
	} else { // server mode
		serverMode(*ports, logger)
	}

	if err != nil {
		logger.Fatal(err)
	}
}

func extractRangeFromParam(param string) (*[]string, error) {
	if !strings.Contains(param, "-") {
		return nil, errors.New("Error: range dont contain a '-'")
	}
	splitted := strings.Split(param, "-")

	if len(splitted) != 2 {
		return nil, errors.New("Error: range is not in format '[num]-[num]'")
	}
	lower, err := strconv.Atoi(splitted[0])
	if err != nil {
		return nil, errors.New("Error: lower range part is not a number")
	}
	upper, err := strconv.Atoi(splitted[1])
	if err != nil {
		return nil, errors.New("Error: upper range part is not a number")
	}
	if lower > upper {
		return nil, errors.New("Error: lower range part is greater then upper range part")
	}
	var ports []string
	for ; lower <= upper; lower++ {
		ports = append(ports, strconv.Itoa(lower))
	}
	return &ports, nil
}

func serverMode(ports []string, logger *log.Logger) {
	logger.Println("running in server mode")

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var servers []*http.Server

	for i, port := range ports {
		addr := ":" + port
		server := &http.Server{Addr: addr, Handler: &handler{i + 1, port, logger}}
		servers = append(servers, server)
		name, err := os.Hostname()
		if err != nil {
			name = "0.0.0.0"
		}

		go func(num int, name, addr string, logger *log.Logger) {
			logger.Printf("Server %d listening on http://%s%s", num, name, addr)

			if err := server.ListenAndServe(); err != nil {
				logger.Fatal(err)
			}
		}(i+1, name, addr, logger)
	}

	<-stop

	logger.Println("\nShutting down the servers...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	for _, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("Error shutting down server: %v\n", err)
		}
	}

	logger.Println("Servers are stopped")
}

type handler struct {
	num    int
	port   string
	logger *log.Logger
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Printf("Recieving connection for server %d on port %s\n", s.num, s.port)
	w.Write([]byte("Ok"))
}

func clientMode(host string, ports []string, logger *log.Logger) {
	logger.Println("running in client mode")

	var wg sync.WaitGroup

	for _, port := range ports {

		wg.Add(1)

		go func(host, port string, logger *log.Logger) {
			defer wg.Done()
			url := fmt.Sprintf("http://%s:%s/", host, port)
			_, err := http.Get(url)
			if err != nil {
				logger.Printf("Connection to %s: Failed\n", url)
			} else {
				logger.Printf("Connection to %s: Success\n", url)
			}
		}(host, port, logger)

	}

	wg.Wait()
}
