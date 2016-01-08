package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/keimoon/gore"
	"github.com/projecteru/eru-agent/logs"
	"github.com/projecteru/eru-ssh/common"
	"github.com/projecteru/eru-ssh/g"
	"golang.org/x/crypto/ssh"
)

func GetRealUserRemote(username string) (string, string) {
	var remote string
	rds := g.GetRedisConn()
	defer g.ReleaseRedisConn(rds)

	keys := strings.Split(username, "~")
	user, ident := keys[0], keys[1]
	routeKey := fmt.Sprintf(common.ROUTE_KEY, user)
	var err error
	var rep *gore.Reply
	if rep, err = gore.NewCommand("HGET", routeKey, ident).Run(rds); err != nil {
		logs.Info("Get info failed", err)
		return "", ""
	}
	if rep.IsNil() {
		return "", ""
	}
	remote, _ = rep.String()
	return user, remote
}

func GetFingerPrint(keyBytes []byte) string {
	h := md5.New()
	h.Write(keyBytes)
	return strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
}

func CheckKey(user, keyHex string) bool {
	rds := g.GetRedisConn()
	defer g.ReleaseRedisConn(rds)
	checkKey := fmt.Sprintf(common.CHECK_KEY, keyHex)
	var err error
	var rep *gore.Reply
	if rep, err = gore.NewCommand("GET", checkKey).Run(rds); err != nil {
		logs.Info("Get info failed", err)
		return false
	}
	if rep.IsNil() {
		return false
	}
	info, _ := rep.String()
	return info == user
}

func LoadKey(keyPath string) (ssh.Signer, error) {
	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	logs.Debug(keyPath, GetFingerPrint(bytes))
	key, err := ssh.ParsePrivateKey(bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
