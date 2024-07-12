package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/koron-go/atermsearch"
)

type Scanner struct {
	verbose       bool
	deviceTimeout time.Duration
}

func (s Scanner) Scan(ctx context.Context) (<-chan *atermsearch.Device, error) {
	ch := make(chan *atermsearch.Device)
	go s.scanAsync(ctx, ch)
	return ch, nil
}

func (s Scanner) scanAsync(ctx context.Context, ch chan<- *atermsearch.Device) {
	var wg sync.WaitGroup
	for i := 0; i < 256; i++ {
		wg.Add(1)
		addr := fmt.Sprintf("192.168.1.%d", i)
		go func() {
			s.search(ctx, ch, addr)
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)
}

func (s Scanner) search(ctx0 context.Context, ch chan<- *atermsearch.Device, addr string) {
	ctx, cancel := context.WithTimeout(ctx0, s.deviceTimeout)
	defer cancel()
	dev, err := atermsearch.Search(ctx, addr)
	if err != nil {
		if s.verbose {
			log.Printf("[INFO] %s: %s", addr, err)
		}
		return
	}
	ch <- dev
}

func main() {
	scanner := Scanner{
		verbose:       false,
		deviceTimeout: 5 * time.Second,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	devices, err := scanner.Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}

	const format = "%-15s\t%-15s\t%-24s\n"
	fmt.Printf(format, "Address", "Product Name", "Mode")
	for dev := range devices {
		fmt.Printf(format, dev.Address, dev.ProductName, dev.SystemMode.Name)
	}
}
