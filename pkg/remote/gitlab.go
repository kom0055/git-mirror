package remote

import (
	"context"
	"fmt"
	"github.com/kom0055/go-flinx"
	_ "github.com/kom0055/go-flinx"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/time/rate"
	"net/http"
	"strings"

	"github.com/kom0055/git-mirror/pkg/utils"
)

var (
	sortAsc         = "asc"
	sortDesc        = "desc"
	orderById       = "id"
	orderByActivity = "last_activity_at"
)

func newGlImpl(ctx context.Context, glUrl, token, groupName, proto string) (Remote, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(glUrl), gitlab.WithCustomLimiter(rate.NewLimiter(5, 3)))
	if err != nil {
		return nil, err
	}
	group, _, err := client.Groups.GetGroup(groupName, &gitlab.GetGroupOptions{}, gitlab.WithContext(ctx))
	if err != nil || group == nil {
		return nil, fmt.Errorf("failed to get group %s: %v", groupName, err)
	}
	return &glRemote{
		proto:  proto,
		group:  group,
		client: client,
	}, nil
}

type glRemote struct {
	proto  string
	group  *gitlab.Group
	client *gitlab.Client
}

func (r *glRemote) GetProjectUrl(ctx context.Context, project *Project) (string, error) {
	group := r.group

	repoName := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s-%s", project.Namespace, project.Name)), " ", "-")
	repoPath := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", group.Name, repoName)), " ", "-")
	repo, resp, err := r.client.Projects.GetProject(repoPath, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		if resp == nil || resp.StatusCode == http.StatusNotFound {
			return "", err

		}
		repo, _, err = r.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        &repoName,
			NamespaceID: &group.ID,

			Visibility: gitlab.Visibility(gitlab.PrivateVisibility),
		}, gitlab.WithContext(ctx))
		if err != nil {
			return "", err
		}
	}

	switch r.proto {
	case utils.HttpProto, utils.HttpsProto:
		return repo.HTTPURLToRepo, nil
	case utils.SshProto, utils.GitProto:

		return repo.SSHURLToRepo, nil

	}
	return "", fmt.Errorf("unknown proto: %s", r.proto)

}

func (r *glRemote) FetchAllProjects(ctx context.Context) ([]*Project, error) {

	allGitlabProjects := []*gitlab.Project{}
	perPage := 50
	for i := 0; ; i++ {
		projects, _, err := r.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    i,
				PerPage: perPage,
			},
			Sort:    &sortAsc,
			OrderBy: &orderById,
		}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if len(projects) == 0 {
			break
		}
		allGitlabProjects = append(allGitlabProjects, projects...)
	}

	allProjects := flinx.ToSlice(flinx.Select(r.fromGitlabProject)(flinx.DistinctBy(func(t *gitlab.Project) int {
		return t.ID
	})(flinx.FromSlice(allGitlabProjects))))

	return allProjects, nil
}

func (r *glRemote) fromGitlabProject(repo *gitlab.Project) *Project {
	var (
		repoUrl string
	)
	switch r.proto {
	case utils.HttpProto, utils.HttpsProto:
		repoUrl = repo.HTTPURLToRepo
	case utils.SshProto, utils.GitProto:

		repoUrl = repo.SSHURLToRepo

	}

	return &Project{
		Name:         repo.Name,
		Namespace:    repo.Namespace.Path,
		URL:          repoUrl,
		RelativePath: repo.PathWithNamespace,
	}
}
