package git

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ip812/ecr-push-notifier/logger"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Target struct {
	RepositroyURL string
	FilePath      string
	Branch        string
	ImageName     string
	ImageTag      string
}

type Git struct {
	dir      string
	repo     *git.Repository
	worktree *git.Worktree
	auth     *http.BasicAuth
	target   Target
	log      logger.Logger
}

func New(log logger.Logger, username, accessToken string, target Target) (*Git, error) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		return nil, err
	}

	auth := &http.BasicAuth{
		Username: username,
		Password: accessToken,
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           target.RepositroyURL,
		Auth:          auth,
		ReferenceName: plumbing.NewBranchReferenceName(target.Branch),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return &Git{
		dir:      dir,
		repo:     repo,
		worktree: worktree,
		auth:     auth,
		target:   target,
		log:      log,
	}, nil
}

func (g *Git) ReplaceImageTag() error {
	fullPath := filepath.Join(g.dir, g.target.FilePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cnt, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(
		fmt.Sprintf(`%s:\d+\.\d+\.\d+`, g.target.ImageName),
	)
	updatedContent := re.ReplaceAll(
		cnt,
		[]byte(fmt.Sprintf("%s:%s", g.target.ImageName, g.target.ImageTag)),
	)

	if err := os.WriteFile(fullPath, []byte(updatedContent), 0644); err != nil {
		return err
	}

	_, err = g.worktree.Add(g.target.FilePath)
	if err != nil {
		return err
	}

	return nil
}

func (g *Git) Push() error {
	_, err := g.worktree.Commit("Update Docker image tag", &git.CommitOptions{
		Author: &object.Signature{
			Name:  fmt.Sprintf("cicd: %s upgraded to %s", g.target.ImageName, g.target.ImageTag),
			Email: "bot.imageupdater@ip812.com",
			When:  time.Now().UTC(),
		},
	})
	if err != nil {
		return err
	}

	err = g.repo.Push(&git.PushOptions{
		Auth: g.auth,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Git) Close() error {
	if err := os.RemoveAll(g.dir); err != nil {
		return err
	}
	return nil
}
