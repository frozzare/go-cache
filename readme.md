# cache [![Build Status](https://travis-ci.org/frozzare/go-cache.svg?branch=master)](https://travis-ci.org/frozzare/go-cache) [![GoDoc](https://godoc.org/github.com/frozzare/go-cache?status.svg)](https://godoc.org/github.com/frozzare/go-cache) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/go-cache)](https://goreportcard.com/report/github.com/frozzare/go-cache)

Go package for dealing with caching. Maybe not so fast.

Requires Go 1.9+ since the package is using [type aliases](https://golang.org/doc/go1.9#language).

## Installation

```
$ go get -u github.com/frozzare/go-cache
```

## Stores

* Memory
* Redis
* Bolt

More cache stores can be implemented by using the provided store interface.

## Example

```go
package main

import (
	"log"

	"github.com/frozzare/go-cache"
	"github.com/frozzare/go-cache/store/redis"
)

func main() {
	c := cache.New(redis.NewStore(&redis.Options{
		Addr: "localhost:6379",
	}))

	if err := c.Set("name", "go"); err != nil {
		log.Fatal(err)
	}

	v, err := c.Get("name")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(v)
}
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)