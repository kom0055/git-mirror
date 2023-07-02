package mirror

import (
	"context"
	"errors"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/kom0055/git-mirror/pkg/utils"
)

type Pusher interface {
	Push(ctx context.Context, repo *git.Repository, repoUrl string) error
}

func NewPusher(permFile, permFilePasswd, userName, token string, proto string) (Pusher, error) {
	auth, err := buildAuth(permFile, permFilePasswd, userName, token, proto)
	if err != nil {
		return nil, err
	}
	return &pusher{
		auth: auth,
	}, nil
}

type pusher struct {
	auth transport.AuthMethod
}

func (s *pusher) Push(ctx context.Context, repo *git.Repository, url string) error {

	remote, err := repo.CreateRemote(&gitconfig.RemoteConfig{
		Name:   utils.DestRemoteName,
		URLs:   []string{url},
		Mirror: utils.IsMirror,
	})
	if err != nil {
		return err
	}

	if err := remote.PushContext(ctx, &git.PushOptions{
		RefSpecs:   []gitconfig.RefSpec{"+refs/tags/*:refs/tags/*", "+refs/heads/*:refs/heads/*", "+refs/merge-requests/*:refs/merge-requests/*"},
		RemoteName: utils.DestRemoteName,
		Auth:       s.auth,
		Force:      true,
		Atomic:     true,
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil

}
