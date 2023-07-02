package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

const (
	IsBare         = true
	IsMirror       = true
	GitUserName    = "git"
	DestRemoteName = "dest"

	GitProto   = "git"
	SshProto   = "ssh"
	HttpProto  = "http"
	HttpsProto = "https"

	TmpPathPattern = "repo"
)

var (
	DefaultTmpPath = fmt.Sprintf("%s/%s/%s", os.TempDir(), TmpPathPattern, uuid.NewString())
)
