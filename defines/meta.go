package defines

import "golang.org/x/crypto/ssh"

type Meta struct {
	Username string
	Password string
	Remote   string
	Client   *ssh.Client
}
