package main

import (
	"fmt"
	"io"

	"github.com/dutchcoders/sshproxy"
	"github.com/projecteru/eru-ssh/g"
	"github.com/projecteru/eru-ssh/proxy"
	"golang.org/x/crypto/ssh"
)

func main() {
	g.LoadConfig()
	g.InitialConn()
	svrConfig := proxy.InitSSHConfig()

	sshproxy.ListenAndServe(
		g.Config.Bind, svrConfig,
		func(conn ssh.ConnMetadata) (*ssh.Client, error) {
			proxy.Lock.RLock()
			defer proxy.Lock.RUnlock()
			meta := proxy.MetaData[conn.RemoteAddr()]
			fmt.Println(meta)
			fmt.Printf("Connection accepted from: %s", conn.RemoteAddr())
			return meta.Client, nil
		},
		func(conn ssh.ConnMetadata, r io.ReadCloser) (io.ReadCloser, error) {
			return sshproxy.NewTypeWriterReadCloser(r), nil
		},
		func(conn ssh.ConnMetadata) error {
			proxy.Lock.Lock()
			defer proxy.Lock.Unlock()
			defer delete(proxy.MetaData, conn.RemoteAddr())
			fmt.Println("Connection closed.")
			return nil
		})
}
