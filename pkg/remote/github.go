package remote

import (
	"context"
	"fmt"
	"github.com/kom0055/go-flinx"
	"golang.org/x/oauth2"
	"net/http"
	"strings"

	"github.com/google/go-github/v47/github"

	"github.com/kom0055/gclone/pkg/utils"
)

var (
	privateVisibility = "private"
)

func newGhImpl(ctx context.Context, token, orgName, proto string) (Remote, error) {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	client := github.NewClient(tc)
	org, _, err := client.Organizations.Get(ctx, orgName)
	if err != nil || org == nil {
		return nil, fmt.Errorf("failed to get org %s: %v", orgName, err)
	}
	return &ghRemote{
		proto:  proto,
		org:    org,
		client: client,
	}, nil
}

type ghRemote struct {
	proto  string
	org    *github.Organization
	client *github.Client
}

func (r *ghRemote) FetchAllProjects(ctx context.Context) ([]*Project, error) {

	allGithubProjects := []*github.Repository{}
	perPage := 50
	for i := 0; ; i++ {
		projects, _, err := r.client.Repositories.List(ctx, "", &github.RepositoryListOptions{
			Sort:      "created",
			Direction: "asc",
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, err
		}
		if len(projects) == 0 {
			break
		}
		allGithubProjects = append(allGithubProjects, projects...)
	}
	allProjects := flinx.ToSlice(flinx.Select(r.fromGithubRepo)(flinx.DistinctBy(func(t *github.Repository) int64 {
		return *t.ID
	})(flinx.FromSlice(allGithubProjects))))

	return allProjects, nil

}

func (r *ghRemote) GetProjectUrl(ctx context.Context, project *Project) (string, error) {
	orgName := *r.org.Login
	repoName := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s-%s", project.Namespace, project.Name)), " ", "-")
	repo, resp, err := r.client.Repositories.Get(ctx, orgName, repoName)
	if err != nil {
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			return "", err

		}
		repo, _, err = r.client.Repositories.Create(ctx, orgName, &github.Repository{
			Name:       &repoName,
			Visibility: &privateVisibility,
		})
		if err != nil {
			return "", err
		}
	}

	switch r.proto {
	case utils.HttpProto, utils.HttpsProto:
		return *repo.CloneURL, nil
	case utils.SshProto, utils.GitProto:
		return *repo.SSHURL, nil

	}
	return "", fmt.Errorf("unknown proto: %s", r.proto)

}

func (r *ghRemote) fromGithubRepo(repo *github.Repository) *Project {
	var (
		repoUrl string
	)
	switch r.proto {
	case utils.HttpProto, utils.HttpsProto:
		repoUrl = *repo.CloneURL
	case utils.SshProto, utils.GitProto:
		repoUrl = *repo.SSHURL

	}
	return &Project{
		Name:         *repo.Name,
		Namespace:    *repo.Owner.Login,
		URL:          repoUrl,
		RelativePath: *repo.FullName,
	}
}
