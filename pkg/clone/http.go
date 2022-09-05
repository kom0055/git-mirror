package clone

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type httpCloner struct {
	auth transport.AuthMethod
	ch   chan struct{}
}

func newHttpCloner(user, token string) (Cloner, error) {
	return &httpCloner{
		auth: &http.BasicAuth{
			Username: user,
			Password: token,
		},
		ch: make(chan struct{}, 4),
	}, nil
}

func (s *httpCloner) Clone(ctx context.Context, localPath, repoUrl string) error {
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
			URL:      repoUrl,
			Auth:     s.auth,
			Tags:     git.AllTags,
			Progress: os.Stdout,
		})
		if err != nil {
			if errors.As(err, &git.ErrRepositoryNotExists) {
				return nil
			}
			return err
		}
	}

	if err = repo.FetchContext(ctx, &git.FetchOptions{
		Auth:     s.auth,
		Tags:     git.AllTags,
		Force:    true,
		Progress: os.Stdout,
	}); err != nil && !errors.As(err, &git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil

}
