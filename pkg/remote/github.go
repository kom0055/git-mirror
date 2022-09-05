package remote

import (
	"context"
	"log"
	"sync"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v47/github"

	"github.com/kom0055/gclone/pkg/utils"
)

type ghRemote struct {
	client *github.Client
}

func newGhImpl(ctx context.Context, token string) (Remote, error) {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	client := github.NewClient(tc)
	return &ghRemote{client: client}, nil
}

func (r *ghRemote) FetchAllProjects(ctx context.Context) []*Project {

	mu := sync.Mutex{}
	allProjects := []*Project{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	utils.GoRoutine(func() {
		defer wg.Done()
		repos, _, err := r.client.Repositories.List(ctx, "", nil)
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		allProjects = append(allProjects, FromGithubRepos(repos...)...)

	})

	wg.Add(1)
	utils.GoRoutine(func() {
		defer wg.Done()
		projects, err := r.lisAllOrgProjects(ctx)
		if err != nil {
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

func (r *ghRemote) lisAllOrgProjects(ctx context.Context) ([]*Project, error) {
	orgs, _, err := r.client.Organizations.List(ctx, "", nil)
	if err != nil {
		return nil, err
	}

	mu := sync.Mutex{}
	allProjects := []*Project{}
	wg := sync.WaitGroup{}
	wg.Add(len(orgs))
	for i := range orgs {
		org := orgs[i]
		utils.GoRoutine(func() {
			defer wg.Done()

			projects, err := r.lisAllTeamProjects(ctx, *org.ID, *org.Name)
			if err != nil {
				log.Println(err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			allProjects = append(allProjects, projects...)
		})

	}
	wg.Wait()

	return allProjects, nil
}

func (r *ghRemote) lisAllTeamProjects(ctx context.Context, orgId int64, orgName string) ([]*Project, error) {
	teams, _, err := r.client.Teams.ListTeams(ctx, orgName, nil)
	if err != nil {
		return nil, err
	}

	mu := sync.Mutex{}
	allProjects := []*Project{}
	wg := sync.WaitGroup{}

	wg.Add(len(teams))
	for i := range teams {
		team := teams[i]
		utils.GoRoutine(func() {
			defer wg.Done()
			repos, _, err := r.client.Teams.ListTeamReposByID(ctx, orgId, *team.ID, nil)
			if err != nil {
				log.Println(err)
				return
			}
			mu.Lock()
			defer mu.Unlock()
			allProjects = append(allProjects, FromGithubRepos(repos...)...)
		})
	}
	wg.Wait()
	return allProjects, nil

}
