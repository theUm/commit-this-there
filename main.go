package main

import (
	"strings"
	"sync"
)

type repo struct {
	Name     string
	FullName string
	URL      string
}

func repoFromUrl(URL string, orgName string) repo {
	return repo{
		Name:     strings.TrimPrefix(URL, "https://"+orgName),
		FullName: strings.TrimPrefix(URL, "https://github.com/"),
		URL:      URL,
	}
}

func main() {
	cfg := &Config{}
	err := ParseEnv(cfg)
	if err != nil {
		panic(err)
	}

	cfg.GithubOrg = "github.com/theUm"
	repos := parseFile(cfg)

	cloneWorkersIn := make(chan repo, 4)
	c := NewCloner(cfg)
	wg := &sync.WaitGroup{}
	for i := 0; i < cfg.CloneWorkersCount; i++ {
		wg.Add(1)
		go c.clone(wg, cloneWorkersIn)
	}

	for _, r := range repos {
		cloneWorkersIn <- r
	}

	close(cloneWorkersIn)

	wg.Wait()
	c.printStats()
}
