package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/koron-go/atermsearch"
)

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
					//log.Printf("[INFO] %s: %s", addr, err)
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
	ctx := context.Background()
	devices, err := scan(ctx, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	const format = "%-15s\t%-15s\t%-24s\n"
	fmt.Printf(format, "Address", "Product Name", "Mode")
	for dev := range devices {
		fmt.Printf(format, dev.Address, dev.ProductName, dev.SystemMode.Name)
	}
}
