package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type githubCloner struct {
	cfg                          Config
	successNum, failNum, ignored *uint64
	rememberFailFunc             func(line string)
}

func NewCloner(cfg *Config) (*githubCloner, error) {
	// todo: deal with references and counters in more civilized way
	successNum := uint64(0)
	failNum := uint64(0)
	ignored := uint64(0)

	rememberFailFunc := func(line string) {}
	if cfg.RememberFails {
		// @todo: make it more concurrency-compatible. now lines may mess-up
		file, err := os.Create(cfg.RememberFailsFilePath)
		if err != nil {
			return nil, fmt.Errorf("create file to write failed repos: %w", err)
		}
		// @todo: close file :)
		//defer file.Close()

		w := bufio.NewWriter(file)

		rememberFailFunc = func(line string) {
			_, _ = fmt.Fprintln(w, line)
			_ = w.Flush()
		}
	}

	return &githubCloner{
		cfg:              *cfg,
		successNum:       &successNum,
		failNum:          &failNum,
		ignored:          &ignored,
		rememberFailFunc: rememberFailFunc,
	}, nil
}

func (g *githubCloner) clone(wg *sync.WaitGroup, in <-chan repo) {
	for repo := range in {
		cloneToDir := filepath.Join(g.cfg.ReposDir, repo.Name)
		dirExists, err := checkDir(cloneToDir)
		if err != nil {
			g.failAndRemember(repo, fmt.Errorf("trying to read dir: %w", err))
			continue
		}
		if dirExists {
			log.Println(fmt.Errorf("[%-30s] repo already exists -- skipping", repo.Name))
			atomic.AddUint64(g.ignored, 1)
			continue
		}
		repoURL := fmt.Sprintf("https://%s:%s@%s/%s", g.cfg.GithubUser, g.cfg.GithubToken, g.cfg.GithubOrg, repo.Name)
		output, err := exec.Command("git", "clone", repoURL, cloneToDir).CombinedOutput()
		if err != nil {
			g.failAndRemember(repo, fmt.Errorf("failed to clone: %s", output))
		} else {
			log.Printf("[%-30s] cloned", repo.FullName)

			if !g.cfg.JustClone {
				err = g.addAndCommit(cloneToDir, repo.Name)
				if err != nil {
					g.failAndRemember(repo, fmt.Errorf("add and commit: %w", err))
					continue
				}

				log.Printf("[%-30s] committed file to %s", repo.FullName, g.cfg.BranchName)

				if g.cfg.DoGitCheck {
					if err = gitCheck(cloneToDir, g.cfg.FileToCommitDestPath); err != nil {
						g.failAndRemember(repo, fmt.Errorf("git check: %w", err))
						continue
					}
				}

				if g.cfg.DoGitPush {
					if err = gitPush(cloneToDir, g.cfg.BranchName); err != nil {
						g.failAndRemember(repo, fmt.Errorf("git push: %w", err))
						continue
					}
					log.Printf("[%-30s] pushed. %s/compare/%s?expand=1", repo.FullName, repo.URL, g.cfg.BranchName)

				}
			}

			// DANGER ZONE! removes cloned repo
			if g.cfg.DeleteCloned {
				err = os.RemoveAll(cloneToDir)
				if err != nil {
					panic(err)
				}
			}

			atomic.AddUint64(g.successNum, 1)
		}
	}
	wg.Done()
}

/*
Actual adding file to repo and cloning part:
* copies file,
* makes branch and checks it out,
* makes git add,
* makes git commit with a message
*/
func (g githubCloner) addAndCommit(cloneToDir, _ string) error {
	fileToCommitDestPath := filepath.Join(cloneToDir, g.cfg.FileToCommitDestPath)

	err := copyFile(g.cfg.FileToCommitSourcePath, fileToCommitDestPath)
	if err != nil {
		return fmt.Errorf("failed to copy file to \"%s\": %w", fileToCommitDestPath, err)
	}

	output, err := exec.Command("git", "-C", cloneToDir, "checkout", "-b", g.cfg.BranchName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to \"git -C %s checkout -b \" to %s: %s", cloneToDir, g.cfg.BranchName, output)
	}

	output, err = exec.Command("git", "-C", cloneToDir, "add", g.cfg.FileToCommitDestPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to \"git -C %s add\" file: %s", cloneToDir, output)
	}

	output, err = exec.Command("git", "-C", cloneToDir, "commit", "-m", g.cfg.CommitMsg).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to \"git -C %s commit\": %w", cloneToDir, output)
	}

	return nil
}

func (g *githubCloner) printStats() {
	log.Printf("all work done: success = %d; ignored = %d; failed = %d \n", *g.successNum, *g.ignored, *g.failNum)
}

func (g *githubCloner) failAndRemember(repo repo, err error) {
	log.Println(fmt.Errorf("[%-30s] %w", repo.Name, err))
	atomic.AddUint64(g.failNum, 1)
	g.rememberFailFunc(fmt.Sprintf("%s -- %s", repo.Name, err.Error()))
}
