package options

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"

	"github.com/kom0055/git-mirror/pkg/mirror"
	"github.com/kom0055/git-mirror/pkg/remote"
	"github.com/kom0055/git-mirror/pkg/utils"
)

type Option struct {
	Worker int
	Source BasicOpt
	Dest   BasicOpt
}

type BasicOpt struct {
	EcdsaPemFile       string
	EcdsaPemFilePasswd string
	RemoteGitlabAddr   string
	User               string
	Token              string
	Proto              string
	GroupName          string
}

func (o *Option) Mirror(ctx context.Context) error {
	defer os.RemoveAll(utils.DefaultTmpPath)

	var (
		source = o.Source
		dest   = o.Dest
	)
	cloner, err := mirror.NewCloner(source.EcdsaPemFile, source.EcdsaPemFilePasswd, source.User, source.Token, source.Proto)
	if err != nil {
		return err
	}

	pusher, err := mirror.NewPusher(dest.EcdsaPemFile, dest.EcdsaPemFilePasswd, dest.User, dest.Token, dest.Proto)
	if err != nil {
		return err
	}

	sourceRm, err := remote.NewRemote(ctx, source.RemoteGitlabAddr, source.Token, source.GroupName, source.Proto)
	if err != nil {
		return err
	}

	destRm, err := remote.NewRemote(ctx, dest.RemoteGitlabAddr, dest.Token, dest.GroupName, dest.Proto)
	if err != nil {
		return err
	}

	projects, err := sourceRm.FetchAllProjects(ctx)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1000)
	worker := o.Worker
	if worker < 1 {
		worker = 1
	}
	ctrl := make(chan struct{}, worker)
	utils.GoRoutine(func() {
		defer close(errCh)
		defer close(ctrl)
		wg := &sync.WaitGroup{}
		finished := &atomic.Int64{}

		mapMutex := &sync.RWMutex{}
		processingProjects := map[string]struct{}{}
		achieves := atomic.Int64{}
		total := int64(len(projects))
		for i := range projects {
			wg.Add(1)
			ctrl <- struct{}{}
			project := projects[i]
			id := project.Identity()
			func() {
				mapMutex.Lock()
				defer mapMutex.Unlock()
				processingProjects[id] = struct{}{}
			}()
			utils.GoRoutine(func() {
				defer func() {
					delete(processingProjects, id)
					processingProjectArr := make([]string, 0, len(processingProjects))
					for name := range processingProjects {
						processingProjectArr = append(processingProjectArr, name)
					}
					log.Printf("total: %v, remain: %v, processing: %v, %v", total, total-achieves.Add(1), len(processingProjects), processingProjectArr)
				}()
				defer func() {
					<-ctrl
				}()
				defer wg.Done()
				defer finished.Add(1)

				ep, err := transport.NewEndpoint(project.URL)
				if err != nil {
					errCh <- fmt.Errorf("parse ssh url %s failed: %s", project.URL, err)
					return
				}
				localPath := filepath.Join(utils.DefaultTmpPath, ep.Host, strings.TrimSuffix(ep.Path, ".git"))
				defer func() {
					_ = os.RemoveAll(localPath)
				}()

				_ = os.RemoveAll(localPath)
				log.Printf("clone %s/%s ", project.Namespace, project.Name)
				repo, err := cloner.Clone(ctx, localPath, project.URL)
				if err != nil {
					if errors.Is(err, transport.ErrEmptyRemoteRepository) {
						return
					}
					log.Printf("clone %s/%s failed: %s", project.Namespace, project.Name, err)
					errCh <- fmt.Errorf("clone %s/%s failed: %s", project.Namespace, project.Name, err)
					return
				}

				destRepoUrl, err := destRm.GetProjectUrl(ctx, project)
				if err != nil {
					log.Printf("get dest repo %s/%s url failed: %v", project.Namespace, project.Name, err)
					errCh <- fmt.Errorf("get dest repo %s/%s url failed: %v", project.Namespace, project.Name, err)
					return
				}

				log.Printf("push %v/%v to %s", project.Namespace, project.Name, destRepoUrl)
				if err := pusher.Push(ctx, repo, destRepoUrl); err != nil {
					log.Printf("push %s failed: %s", destRepoUrl, err)
					errCh <- fmt.Errorf("push %s failed: %s", destRepoUrl, err)
					return
				}
			})
		}

		wg.Wait()
	})

	var multiErr *multierror.Error
	for err := range errCh {
		multiErr = multierror.Append(multiErr, err)
	}

	return multiErr.ErrorOrNil()
}
