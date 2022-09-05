package options

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"gopkg.in/yaml.v3"

	"github.com/kom0055/gclone/pkg/clone"
	"github.com/kom0055/gclone/pkg/config"
	"github.com/kom0055/gclone/pkg/remote"
	"github.com/kom0055/gclone/pkg/utils"
)

type Option struct {
	CfgFilePath        string
	RepoRootPath       string
	EcdsaPemFile       string
	EcdsaPemFilePasswd string
	RemoteGitlabAddr   string
	SyncFromRemote     bool

	User  string
	Token string
}

func (o *Option) Complete() error {

	cfg := config.Cfg{}
	if len(o.CfgFilePath) == 0 {
		o.CfgFilePath = filepath.Join(config.HomeDir, config.CfgFilePath)
	}
	if len(o.CfgFilePath) == 0 {
		return nil
	}
	b, err := os.ReadFile(o.CfgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	if len(o.RepoRootPath) == 0 {
		o.RepoRootPath = cfg.RepoRootPath
	}

	if len(o.EcdsaPemFile) == 0 {
		o.EcdsaPemFile = cfg.EcdsaPemFile
	}

	if len(o.EcdsaPemFilePasswd) == 0 {
		o.EcdsaPemFilePasswd = cfg.EcdsaPemFilePasswd
	}

	if len(o.User) == 0 {
		o.User = cfg.User
	}
	if len(o.Token) == 0 {
		o.Token = cfg.Token
	}

	if len(o.RemoteGitlabAddr) == 0 {
		o.RemoteGitlabAddr = cfg.RemoteGitlabAddr
	}

	return nil

}

func (o *Option) Run(ctx context.Context, args []string) error {
	if !o.SyncFromRemote {
		repoUrl := args[0]
		return o.CloneRepo(ctx, repoUrl)
	}
	if len(args) > 0 {
		return fmt.Errorf("unsupport args %v", args)
	}
	return o.SyncRemote(ctx)

}

func (o *Option) CloneRepo(ctx context.Context, repoUrl string) error {
	ep, err := transport.NewEndpoint(repoUrl)
	if err != nil {
		return err
	}
	localPath := filepath.Join(o.RepoRootPath, ep.Host, strings.TrimSuffix(ep.Path, ".git"))
	cloner, err := clone.NewCloner(o.EcdsaPemFile, o.EcdsaPemFilePasswd, o.User, o.Token, ep.Protocol)
	if err != nil {
		return err
	}

	if err = cloner.Clone(ctx, localPath, repoUrl); err != nil {
		return err
	}
	return nil
}

func (o *Option) SyncRemote(ctx context.Context) error {

	cloner, err := clone.NewCloner(o.EcdsaPemFile, o.EcdsaPemFilePasswd, o.User, o.Token, config.SshProto)
	if err != nil {
		return err
	}
	rm, err := remote.NewRemote(ctx, o.RemoteGitlabAddr, o.Token)
	if err != nil {
		return err
	}
	projects := rm.FetchAllProjects(ctx)
	wg := sync.WaitGroup{}
	wg.Add(len(projects))

	for i := range projects {
		project := projects[i]
		utils.GoRoutine(func() {
			defer wg.Done()
			ep, err := transport.NewEndpoint(project.SSHUrl)
			if err != nil {
				log.Println(err)
				return
			}
			localPath := filepath.Join(o.RepoRootPath, ep.Host, strings.TrimSuffix(ep.Path, ".git"))
			if err = cloner.Clone(ctx, localPath, project.SSHUrl); err != nil {
				log.Println(err)
				return
			}

		})

	}
	wg.Wait()
	return nil
}
