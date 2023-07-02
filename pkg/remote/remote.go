package remote

import (
	"context"
)

type Remote interface {
	FetchAllProjects(ctx context.Context) ([]*Project, error)
	GetProjectUrl(ctx context.Context, project *Project) (string, error)
}

func NewRemote(ctx context.Context, glUrl, token, org, proto string) (Remote, error) {
	if len(glUrl) == 0 {
		return newGhImpl(ctx, token, org, proto)
	}

	return newGlImpl(ctx, glUrl, token, org, proto)
}

type Project struct {
	Name         string
	Namespace    string
	URL          string
	RelativePath string
}
