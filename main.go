package main

import (
	"github.com/projecteru/eru-ssh/g"
	"github.com/projecteru/eru-ssh/proxy"
)

func main() {
	g.LoadConfig()
	g.InitialConn()

	proxy.ListenAndServe()
}
