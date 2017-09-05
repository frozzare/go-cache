# cache [![Build Status](https://travis-ci.org/frozzare/go-cache.svg?branch=master)](https://travis-ci.org/frozzare/go-cache) [![GoDoc](https://godoc.org/github.com/frozzare/go-cache?status.svg)](https://godoc.org/github.com/frozzare/go-cache) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/go-cache)](https://goreportcard.com/report/github.com/frozzare/go-cache)

Go package for dealing with caching. The package has only [redis](https://redis.io/) support right now but more cache stores can be implemented using the store interface.

Requires Go 1.9+ since the package is using [type aliases](https://golang.org/doc/go1.9#language).

## Installation

```
$ go get -u github.com/frozzare/go-cache
```

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/frozzare/go-cache"
	"github.com/frozzare/go-cache/store/redis"
)

func main() {
	c := cache.New(cache.Redis(&redis.Options{
		Addr: "localhost:6379",
	}))

	if err := c.Set("name", "go"); err != nil {
		log.Fatal(err)
	}

	v, err := c.Get("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)
}
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)