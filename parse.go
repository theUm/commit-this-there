package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
)

func parseFile(cfg *Config) []repo {
	var repos []repo

	file, err := os.Open(cfg.ReposListFile)
	if err != nil {
		log.Fatalf("failed to open repos list file: %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		r := repoFromUrl(line, cfg.GithubOrg+"/")
		_, err := url.Parse(line)
		if err != nil {
			log.Println("skipping line -- cant parse url:" + line)
			continue
		}
		repos = append(repos, r)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed scan line: %s", err.Error())
	}
	return repos
}
