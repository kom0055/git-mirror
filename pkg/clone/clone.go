package clone

import (
	"context"
	"fmt"

	"github.com/kom0055/gclone/pkg/config"
)

const (
	isBare      = false
	gitUserName = "git"
)

type Cloner interface {
	Clone(ctx context.Context, localPath, repoUrl string) error
}

func NewCloner(permFile, permFilePasswd, userName, token string, proto string) (Cloner, error) {
	switch proto {
	case config.HttpProto, config.HttpsProto:
		return newHttpCloner(userName, token)
	case config.SshProto, config.GitProto:
		return newSshCloner(gitUserName, permFile, permFilePasswd)
	default:

	}
	return nil, fmt.Errorf("unsupport protocol %s", proto)
}
