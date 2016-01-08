package proxy

import (
	"errors"
	"net"
	"sync"

	"github.com/projecteru/eru-agent/logs"
	"github.com/projecteru/eru-ssh/defines"
	"github.com/projecteru/eru-ssh/g"
	"github.com/projecteru/eru-ssh/utils"
	"golang.org/x/crypto/ssh"
)

var Lock sync.RWMutex
var MetaData map[net.Addr]defines.Meta = map[net.Addr]defines.Meta{}

func InitSSHConfig() *ssh.ServerConfig {
	privKey, err := utils.LoadKey(g.Config.PrivKey)
	if err != nil {
		logs.Assert(err, "Failed to load priv key")
	}

	config := &ssh.ServerConfig{
		AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
			logs.Debug("Method type", method, "Error", err)
		},
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			username := conn.User()
			clientAddr := conn.RemoteAddr()
			keyHex := utils.GetFingerPrint(key.Marshal())

			user, remote := utils.GetRealUserRemote(username)
			if user == "" || remote == "" {
				return nil, errors.New("Wrong info")
			}

			logs.Info("Login attempt", conn.RemoteAddr(), username, user, remote, keyHex)

			if !utils.CheckKey(user, keyHex) {
				return nil, errors.New("Wrong key")
			}

			meta := defines.Meta{
				Username: username,
				Pubkey:   key,
			}

			clientConfig := &ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(privKey),
				},
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
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			username := conn.User()
			clientAddr := conn.RemoteAddr()

			user, remote := utils.GetRealUserRemote(username)
			if user == "" || remote == "" {
				return nil, errors.New("Wrong info")
			}

			logs.Info("Login attempt", conn.RemoteAddr(), username, string(password))

			meta := defines.Meta{
				Username: username,
				Password: string(password),
			}

			clientConfig := &ssh.ClientConfig{
				User: "root",
				Auth: []ssh.AuthMethod{
					ssh.Password(string(password)),
				},
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

	hostKey, err := utils.LoadKey(g.Config.HostKey)
	if err != nil {
		logs.Assert(err, "Failed to load host key")
	}
	config.AddHostKey(hostKey)
	return config
}
