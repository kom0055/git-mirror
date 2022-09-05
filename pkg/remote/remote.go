package remote

import (
	"context"

	"github.com/google/go-github/v47/github"
	"github.com/xanzy/go-gitlab"
)

type Remote interface {
	FetchAllProjects(ctx context.Context) []*Project
}

func NewRemote(ctx context.Context, glUrl, token string) (Remote, error) {
	if len(glUrl) == 0 {
		return newGhImpl(ctx, token)
	}

	return newGlImpl(glUrl, token)
}

type Project struct {
	Name         string
	SSHUrl       string
	RelativePath string
}

func FromGitlabProjects(glProjects ...*gitlab.Project) []*Project {
	projects := make([]*Project, len(glProjects))
	for i := range glProjects {
		projects[i] = &Project{
			Name:         glProjects[i].Name,
			SSHUrl:       glProjects[i].SSHURLToRepo,
			RelativePath: glProjects[i].PathWithNamespace,
		}
	}
	return projects
}

func FromGithubRepos(ghRepos ...*github.Repository) []*Project {
	projects := make([]*Project, len(ghRepos))
	for i := range ghRepos {
		projects[i] = &Project{
			Name:         *ghRepos[i].Name,
			SSHUrl:       *ghRepos[i].SSHURL,
			RelativePath: *ghRepos[i].FullName,
		}
	}
	return projects
}
