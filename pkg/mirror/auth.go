package mirror

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/kom0055/git-mirror/pkg/utils"
)

func buildAuth(permFile, permFilePasswd, userName, token string, proto string) (transport.AuthMethod, error) {
	switch proto {
	case utils.HttpProto, utils.HttpsProto:
		return &http.BasicAuth{
			Username: userName,
			Password: token,
		}, nil
	case utils.SshProto, utils.GitProto:
		publicKeys, err := ssh.NewPublicKeysFromFile(utils.GitUserName, permFile, permFilePasswd)
		if err != nil {
			return nil, err
		}
		publicKeys.HostKeyCallback = utils.IgnoreHostKeyCB
		return publicKeys, nil
	}
	return nil, fmt.Errorf("unsupport protocol %s", proto)
}
