package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/koron-go/atermsearch"
)

var ignorable = []string{
	"no route to host",
	"connection refused",
	"network is unreachable",
	"failed to HTTP request with status",
}

func isIgnorable(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	msg := err.Error()
	for _, ignore := range ignorable {
		if strings.Contains(msg, ignore) {
			return true
		}
	}
	return false
}

func scan(ctx context.Context, deviceTimeout time.Duration) (<-chan *atermsearch.Device, error) {
	ch := make(chan *atermsearch.Device)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 256; i++ {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				ctx2, cancel := context.WithTimeout(ctx, deviceTimeout)
				defer cancel()
				dev, err := atermsearch.Search(ctx2, addr)
				if err != nil {
					if !isIgnorable(err) {
						slog.Warn("error cannot be ignored", "address", addr, "error", err)
					}
					return
				}
				ch <- dev
			}(fmt.Sprintf("192.168.1.%d", i))
		}
		wg.Wait()
		close(ch)
	}()
	return ch, nil
}

func main() {
	var (
		timeout time.Duration
	)
	flag.DurationVar(&timeout, "timeout", 8*time.Second, "timeout for each address")
	flag.Parse()

	ctx := context.Background()
	devices, err := scan(ctx, timeout)
	if err != nil {
		log.Fatal(err)
	}
	const format = "%-22s\t%-22s\t%-24s\n"
	fmt.Printf(format, "Address", "Product Name", "Mode")
	for dev := range devices {
		fmt.Printf(format, dev.Address, dev.ProductName, dev.SystemMode.Name)
	}
}
