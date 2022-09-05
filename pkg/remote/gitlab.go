package remote

import (
	"context"
	"log"
	"sync"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/time/rate"

	"github.com/kom0055/gclone/pkg/utils"
)

type glRemote struct {
	client *gitlab.Client
}

func newGlImpl(glUrl, token string) (Remote, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(glUrl), gitlab.WithCustomLimiter(rate.NewLimiter(5, 3)))
	if err != nil {
		return nil, err
	}
	return &glRemote{client: client}, nil
}

func (r *glRemote) FetchAllProjects(ctx context.Context) []*Project {

	mu := sync.Mutex{}
	allProjects := []*Project{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	utils.GoRoutine(func() {
		defer wg.Done()
		projects, _, err := r.client.Projects.ListProjects(&gitlab.ListProjectsOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		allProjects = append(allProjects, FromGitlabProjects(projects...)...)

	})

	wg.Add(1)
	utils.GoRoutine(func() {
		defer wg.Done()
		projects, err := r.listAllGroupProjects(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		allProjects = append(allProjects, projects...)

	})
	wg.Wait()
	set := utils.NewAnySet[string]()
	res := []*Project{}
	for i := range allProjects {
		project := allProjects[i]
		if set.Has(project.SSHUrl) {
			continue
		}
		set.Insert(project.SSHUrl)
		res = append(res, project)
	}
	return res
}

func (r *glRemote) listAllGroupProjects(ctx context.Context) ([]*Project, error) {
	mu := sync.Mutex{}
	allProjects := []*Project{}
	groups, _, err := r.client.Groups.ListGroups(&gitlab.ListGroupsOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(groups))
	for i := range groups {
		group := groups[i]
		utils.GoRoutine(func() {
			defer wg.Done()
			projects, _, err := r.client.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{}, gitlab.WithContext(ctx))
			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			allProjects = append(allProjects, FromGitlabProjects(projects...)...)
		})
	}
	wg.Wait()
	return allProjects, nil
}
