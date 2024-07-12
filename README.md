# koron-go/atermsearch

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/atermsearch)](https://pkg.go.dev/github.com/koron-go/atermsearch)
[![Actions/Go](https://github.com/koron-go/atermsearch/workflows/Go/badge.svg)](https://github.com/koron-go/atermsearch/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/atermsearch)](https://goreportcard.com/report/github.com/koron-go/atermsearch)

Package to search for [Aterm devices](https://www.aterm.jp/product/atermstation/).
It is reimplementation of [Aterm search tool](https://www.aterm.jp/web/model/aterm_search.html) in Go.

## Getting Started

``` console
$ go get github.com/koron-go/atermsearch@latest
```


Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/koron-go/atermsearch"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    dev, err := atermsearch.Search(ctx, "192.168.0.123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Aterm device: %+v\n", dev)
}
```
