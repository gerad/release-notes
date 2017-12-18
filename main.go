package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Config struct {
	base     string
	head     string
	owner    string
	password string
	repo     string
	username string
}

type Github struct {
	username string
	password string
}

func main() {
	config := parseFlags()
	github := Github{username: config.username, password: config.password}

	var repo struct {
		Commits []struct {
			Commit struct {
				Author struct {
					Name string
				}
				Message string
			}
		}
	}
	err := github.Get(fmt.Sprintf("repos/%v/%v/compare/%v...%v",
		config.owner,
		config.repo,
		config.base,
		config.head), &repo)

	if err != nil {
		log.Fatalf("%v", err)
	}

	re1 := regexp.MustCompile("Merge pull request #(?P<num>\\d+).*\n\n(?P<title>.*)")
	re2 := regexp.MustCompile("(?P<title>.*) \\(#(?P<num>\\d+)\\)")

	for _, commit := range repo.Commits {
		for _, re := range []*regexp.Regexp{re1, re2} {
			matches := re.FindStringSubmatch(commit.Commit.Message)
			if matches != nil && len(matches) == 3 {
				result := make(map[string]string)
				for i, name := range re.SubexpNames() {
					if i != 0 {
						result[name] = matches[i]
					}
				}

				fmt.Printf("%s - https://github.com/%s/%s/pull/%s - %s\n", commit.Commit.Author.Name, config.owner, config.repo, result["num"], result["title"])
				break
			}
		}
	}
}

func (config Github) Get(path string, v interface{}) error {
	// log.Printf("GET %s", path)

	resp, err := http.Get(fmt.Sprintf("https://%v:%v@api.github.com/%v",
		config.username,
		config.password,
		path))

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func parseFlags() (config Config) {

	flag.StringVar(&config.base, "base", "demo", "base branch for comparison")
	flag.StringVar(&config.head, "head", "master", "head branch for comparison")
	flag.StringVar(&config.password, "password", "", "(REQUIRED) github password")
	flag.StringVar(&config.username, "username", "", "(REQUIRED) github username")
	flag.StringVar(&config.owner, "owner", "${username}", "github repo owner")
	flag.StringVar(&config.repo, "repo", "", "(REQUIRED) github repo name")

	flag.Parse()

	if config.username == "" || config.password == "" || config.repo == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if config.owner == "${username}" {
		config.owner = config.username
	}

	return
}
