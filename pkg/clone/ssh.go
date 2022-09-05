package clone

import (
	"context"
	"errors"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type sshCloner struct {
	auth transport.AuthMethod
	ch   chan struct{}
}

func newSshCloner(user, permFile, passwd string) (Cloner, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile(user, permFile, passwd)
	if err != nil {
		return nil, err
	}
	publicKeys.HostKeyCallback = IgnoreHostKeyCB
	return &sshCloner{
		auth: publicKeys,
		ch:   make(chan struct{}, 4),
	}, nil
}

func (s *sshCloner) Clone(ctx context.Context, localPath, repoUrl string) error {
	s.ch <- struct{}{}
	log.Println("cloning", repoUrl, localPath, "start")
	defer func() {
		<-s.ch
		log.Println("cloning", repoUrl, "done")
	}()
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		if !errors.As(err, &git.ErrRepositoryNotExists) {
			return err
		}
		repo, err = git.PlainCloneContext(ctx, localPath, isBare, &git.CloneOptions{
			URL:  repoUrl,
			Auth: s.auth,
			Tags: git.AllTags,
		})

		if err != nil {
			if errors.As(err, &git.ErrRepositoryNotExists) {
				return nil
			}
			return err
		}
	}

	if err = repo.FetchContext(ctx, &git.FetchOptions{
		Auth:  s.auth,
		Tags:  git.AllTags,
		Force: true,
	}); err != nil && !errors.As(err, &git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil

}
