package g

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/projecteru/eru-agent/logs"
	"github.com/projecteru/eru-ssh/common"
	"github.com/projecteru/eru-ssh/defines"
	"gopkg.in/yaml.v2"
)

var Config = defines.SSHConfig{}

func LoadConfig() {
	var configPath string
	var version bool
	flag.BoolVar(&logs.Mode, "DEBUG", false, "enable debug")
	flag.StringVar(&configPath, "c", "ssh.yaml", "config file")
	flag.BoolVar(&version, "v", false, "show version")
	flag.Parse()
	if version {
		logs.Info("Version", common.VERSION)
		os.Exit(0)
	}
	load(configPath)
}

func load(configPath string) {
	if _, err := os.Stat(configPath); err != nil {
		logs.Assert(err, "config file invaild")
	}

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		logs.Assert(err, "Read config file failed")
	}

	if err := yaml.Unmarshal(b, &Config); err != nil {
		logs.Assert(err, "Load config file failed")
	}

	logs.Debug("Configure:", Config)
}
