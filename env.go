package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	GithubOrg              string `env:"GITHUB_ORG,required"`                                    // your Github org name in format of "github.com/orgName"
	GithubUser             string `env:"GITHUB_USER,required"`                                   // your Github username
	GithubToken            string `env:"GITHUB_TOKEN,required"`                                  // your Github token from https://github.com/settings/tokens
	ReposListFile          string `env:"REPOS_LIST_FILE" envDefault:"./repos.txt"`               // list of repos
	ReposDir               string `env:"REPOS_DIR" envDefault:"/tmp/cloned_repos"`               // dir to clone repos to
	CloneWorkersCount      int    `env:"CLONE_WORKERS" envDefault:"4"`                           // num of concurrency for cloners
	JustClone              bool   `env:"JUST_CLONE" envDefault:"false"`                          // just clone repos, no commit no push, nothing
	FileToCommitSourcePath string `env:"FILE_TO_COMMIT_SRC" envDefault:"./dependabot.yml"`       // file to be added to git
	FileToCommitDestPath   string `env:"FILE_TO_COMMIT_DST" envDefault:".github/dependabot.yml"` // destination for file to be cloned inside the repo
	BranchName             string `env:"BRANCH_NAME" envDefault:"add-dependabot"`                // branch name
	CommitMsg              string `env:"COMMIT_MSG" envDefault:"adding dependabot"`              // commit message
	DoGitCheck             bool   `env:"DO_GIT_CHECK" envDefault:"true"`                         // check last commit to contain specified file ( runs before push)
	DoGitPush              bool   `env:"DO_GIT_PUSH" envDefault:"false"`                         // to push or not to push - that is the question
	DeleteCloned           bool   `env:"DELETE_CLONED" envDefault:"false"`                       // delete cloned repo(repos) if everything is fine
}

func ParseEnv(envStruct *Config) error {
	if err := env.Parse(envStruct); err != nil {
		return err
	}
	return nil
}
