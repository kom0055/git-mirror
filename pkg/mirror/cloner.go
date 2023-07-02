package mirror

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/kom0055/git-mirror/pkg/utils"
)

type Cloner interface {
	Clone(ctx context.Context, localPath, repoUrl string) (*git.Repository, error)
}

func NewCloner(permFile, permFilePasswd, userName, token string, proto string) (Cloner, error) {
	auth, err := buildAuth(permFile, permFilePasswd, userName, token, proto)
	if err != nil {
		return nil, err
	}
	return &cloner{
		auth: auth,
	}, nil
}

type cloner struct {
	auth transport.AuthMethod
}

func (s *cloner) Clone(ctx context.Context, localPath, repoUrl string) (*git.Repository, error) {

	repo, err := git.PlainCloneContext(ctx, localPath, utils.IsBare, &git.CloneOptions{
		URL:    repoUrl,
		Auth:   s.auth,
		Tags:   git.AllTags,
		Mirror: utils.IsMirror,
	})

	if err != nil {
		return nil, err
	}
	return repo, nil

}
