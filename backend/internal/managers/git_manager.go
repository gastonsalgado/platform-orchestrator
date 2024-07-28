package managers

import (
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"
)

type gitManager struct {
	BasePath                 string
	Url                      string
	Reference                string
	InfraTenantTemplatesPath string
	InfraTenantsPath         string
	Repository               *git.Repository
	Worktree                 *git.Worktree
}

var gitManagerLock = &sync.Mutex{}
var gitManagerInstance *gitManager
var Logger *zap.Logger

func GetGitManagerInstance() *gitManager {
	if gitManagerInstance == nil {
		gitManagerLock.Lock()
		defer gitManagerLock.Unlock()
		if gitManagerInstance == nil {
			gitManagerInstance = &gitManager{}
		}
	}

	return gitManagerInstance
}

func (g *gitManager) Init() error {
	g.BasePath = os.Getenv("GIT_BASE_PATH")
	g.Url = os.Getenv("GIT_URL")
	g.Reference = os.Getenv("GIT_REFERENCE")
	g.InfraTenantTemplatesPath = os.Getenv("GIT_INFRA_TENANT_TEMPLATES_PATH")
	g.InfraTenantsPath = os.Getenv("GIT_INFRA_TENANTS")

	var repository *git.Repository

	_, err := os.Stat(g.BasePath)
	if os.IsNotExist(err) {
		repository, err = git.PlainClone(g.BasePath, false, &git.CloneOptions{
			URL: g.Url,
		})
	} else {
		repository, err = git.PlainOpen(g.BasePath)
	}
	if err != nil {
		return err
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	g.Repository = repository
	g.Worktree = worktree
	return nil
}

func (g *gitManager) Pull() error {
	Logger.Info("git pull origin")
	err := g.Worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		Logger.Error(err.Error())
	}

	return err
}

func (g *gitManager) Push(path string, commitMessage string) error {
	Logger.Info("git add")
	_, err := g.Worktree.Add(path)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	Logger.Info("git status --porcelain")
	status, err := g.Worktree.Status()
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	Logger.Debug(status.String())

	Logger.Info("git commit")
	// commit, err := w.Commit("example go-git commit", &git.CommitOptions{
	// 	Author: &object.Signature{
	// 		Name:  "John Doe",
	// 		Email: "john@doe.org",
	// 		When:  time.Now(),
	// 	},
	// })
	commit, err := g.Worktree.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	Logger.Info("git show -s")
	obj, err := g.Repository.CommitObject(commit)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	Logger.Debug(obj.Message)

	Logger.Info("git push")
	err = g.Repository.Push(&git.PushOptions{})
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}
