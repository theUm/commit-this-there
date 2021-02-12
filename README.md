Tool to quickly commit single file across multiple repos.

It clones repos via cli call `git clone`, copies file there, and commits and pushes (optionally) it to specified branch.

Next envs are used:
```
"GITHUB_ORG,required"                                    // your Github org name in format of "github.com/orgName"
"GITHUB_USER,required"                                   // your Github username
"GITHUB_TOKEN,required"                                  // your Github token from https://github.com/settings/tokens
"REPOS_LIST_FILE" envDefault:"./repos.txt"               // list of repos
"REPOS_DIR" envDefault:"/tmp/cloned_repos"               // dir to clone repos to
"CLONE_WORKERS" envDefault:"4"                           // num of concurrency for cloners
"JUST_CLONE" envDefault:"false"                          // just clone repos, no commit no push, nothing
"FILE_TO_COMMIT_SRC" envDefault:"./dependabot.yml"       // file to be added to git
"FILE_TO_COMMIT_DST" envDefault:".github/dependabot.yml" // destination for file to be cloned inside the repo
"BRANCH_NAME" envDefault:"add-dependabot"         // branch name
"COMMIT_MSG" envDefault:"adding dependabot"              // commit message
"DO_GIT_CHECK" envDefault:"true"                         // check last commit to contain specified file ( runs before push)
"DO_GIT_PUSH" envDefault:"false"                         // to push or not to push - that is the question
"DELETE_CLONED" envDefault:"false"                       // delete cloned repo(repos) if everything is fine
```

Note: you have to manually delete failed repos if some step fails. It would be like this for now