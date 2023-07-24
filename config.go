package dblite

import "github.com/lesismal/nbio/nbhttp"

var Server = nbhttp.NewServer(nbhttp.Config{Network: "tcp", Addrs: []string{"localhost:1111"}})
