package proxy

import (
	"errors"
	"io/ioutil"
	"net"
	"strings"
	"sync"

	"github.com/keimoon/gore"
	"github.com/projecteru/eru-agent/logs"
	"github.com/projecteru/eru-ssh/defines"
	"github.com/projecteru/eru-ssh/g"
	"golang.org/x/crypto/ssh"
)

var Lock sync.RWMutex
var MetaData map[net.Addr]defines.Meta = map[net.Addr]defines.Meta{}

func InitSSHConfig() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			logs.Info("Login attempt", conn.RemoteAddr(), conn.User(), string(password))
			clientAddr := conn.RemoteAddr()
			meta := defines.Meta{
				Username: conn.User(),
				Password: string(password),
			}

			clientConfig := &ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.Password(string(password)),
				},
			}
			rds := g.GetRedisConn()
			defer g.ReleaseRedisConn(rds)
			var remote string

			keys := strings.Split(conn.User(), "~")
			user, host := keys[0], keys[1]
			if rep, err := gore.NewCommand("HGET", user, host).Run(rds); err != nil {
				return nil, err
			} else {
				if rep.IsNil() {
					return nil, errors.New("no dest")
				}
				remote, _ = rep.String()
			}

			client, err := ssh.Dial("tcp", remote, clientConfig)
			if err != nil {
				return nil, err
			}
			meta.Remote = remote
			meta.Client = client
			Lock.Lock()
			defer Lock.Unlock()
			MetaData[clientAddr] = meta
			return nil, nil
		},
	}
	privateBytes, err := ioutil.ReadFile(g.Config.Key)
	if err != nil {
		logs.Assert(err, "Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		logs.Assert(err, "Failed to parse private key")
	}

	config.AddHostKey(private)
	return config
}
