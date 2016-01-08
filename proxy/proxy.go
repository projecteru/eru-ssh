package proxy

import (
	"io"
	"net"

	"github.com/projecteru/eru-agent/logs"
	"github.com/projecteru/eru-ssh/g"
	"golang.org/x/crypto/ssh"
)

type SSHConn struct {
	net.Conn
	config *ssh.ServerConfig
}

func (self *SSHConn) serve() error {
	serverConn, chans, reqs, err := ssh.NewServerConn(self.Conn, self.config)
	if err != nil {
		return err
	}
	defer serverConn.Close()

	clientConn, err := getClient(serverConn)
	if err != nil {
		return err
	}
	defer clientConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		remoteChannel, remoteRequest, err := clientConn.OpenChannel(newChannel.ChannelType(), newChannel.ExtraData())
		if err != nil {
			return err
		}

		localChannel, localRequest, err := newChannel.Accept()
		if err != nil {
			return err
		}

		// connect requests
		go func() {
			logs.Debug("Waiting for request")
		r:
			for {
				var req *ssh.Request
				var dst ssh.Channel

				select {
				case req = <-localRequest:
					dst = remoteChannel
					logs.Debug("from local to remote")
				case req = <-remoteRequest:
					dst = localChannel
					logs.Debug("from remote to local")
				}

				if req == nil {
					break
				}

				logs.Debug("Request", req.Type, req.WantReply)
				b, err := dst.SendRequest(req.Type, req.WantReply, req.Payload)
				if err != nil {
					logs.Info(err)
				}
				if req.WantReply {
					req.Reply(b, nil)
				}
				switch req.Type {
				case "exit-status":
					break r
				}
			}

			localChannel.Close()
			remoteChannel.Close()
		}()

		// connect channels
		logs.Debug("Connecting channels")

		go io.Copy(remoteChannel, localChannel)
		go io.Copy(localChannel, remoteChannel)

		defer remoteChannel.Close()
		defer localChannel.Close()
	}

	closeConn(serverConn)
	return nil
}

func ListenAndServe() error {
	svrConfig := InitSSHConfig()
	listener, err := net.Listen("tcp", g.Config.Bind)
	if err != nil {
		logs.Assert(err, "net.Listen failed")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Assert(err, "listen.Accept failed")
		}

		sshConn := &SSHConn{conn, svrConfig}

		go func() {
			if err := sshConn.serve(); err != nil {
				logs.Info("Error occured while serving", err)
				return
			}
			logs.Info("Connection closed.")
		}()
	}
}

func getClient(conn ssh.ConnMetadata) (*ssh.Client, error) {
	Lock.RLock()
	defer Lock.RUnlock()
	meta := MetaData[conn.RemoteAddr()]
	logs.Debug("Connection accepted from", conn.RemoteAddr())
	return meta.Client, nil
}

func wrap(conn ssh.ConnMetadata, r io.ReadCloser) (io.ReadCloser, error) {
	return r, nil
}

func closeConn(conn ssh.ConnMetadata) error {
	Lock.Lock()
	defer Lock.Unlock()
	defer delete(MetaData, conn.RemoteAddr())
	logs.Debug("Clean sessions")
	return nil
}
