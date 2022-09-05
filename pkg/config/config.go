package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
)

const (
	CfgFilePath = ".gclone.yaml"

	GitProto   = "git"
	SshProto   = "ssh"
	HttpProto  = "http"
	HttpsProto = "https"
)

var (
	HomeDir string
)

func init() {
	var err error
	HomeDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
}

type Cfg struct {
	RepoRootPath       string `json:"goPath,omitempty" yaml:"goPath,omitempty"`
	EcdsaPemFile       string `json:"ecdsaPemFile,omitempty" yaml:"ecdsaPemFile,omitempty"`
	EcdsaPemFilePasswd string `json:"ecdsaPemFilePasswd" yaml:"ecdsaPemFilePasswd"`
	RemoteGitlabAddr   string `json:"remoteGitlabAddr" yaml:"remoteGitlabAddr"`

	User  string `default:"git" json:"user,omitempty" yaml:"user,omitempty"`
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

type plain Cfg

func (c *Cfg) UnmarshalYAML(unmarshal func(interface{}) error) error {

	if err := defaults.Set(c); err != nil {
		return err
	}
	c.RepoRootPath = filepath.Join(HomeDir, "go/src")
	c.EcdsaPemFile = filepath.Join(HomeDir, ".ssh/id_ecdsa")
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil

}
