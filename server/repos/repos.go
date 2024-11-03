package repos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	utils "github.com/johnietre/utils/go"
)

var (
	repos = utils.NewAValue([]RepoInfo{})
)

func InitRepos() error {
	// Load the GitHub repos
	if err := RefreshRepos(); err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			if err := RefreshRepos(); err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}

func RefreshRepos() error {
	resp, err := http.Get(
		"https://api.github.com/users/johnietre/repos?sort=pushed&direction=desc&per_page=3",
	)
	if err != nil {
		return fmt.Errorf("Error getting repos json: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Received non-200 status: %s", resp.Status)
	}
	reposMaps := []RepoInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&reposMaps); err != nil {
		return fmt.Errorf("Error decoding repos json: %v", err)
	}
	reposArr := make([]RepoInfo, 0, len(reposMaps))
	for _, repo := range reposMaps {
		reposArr = append(reposArr, repo)
		/*
		   iUrl, ok := repo["html_url"]
		   if !ok {
		     logger.Println("Missing repo html_url")
		   }
		   url, ok := iUrl.(string)
		   if !ok {
		     logger.Println("Invalid repo html_url received, got %v", iUrl)
		   }
		   iName, ok := repo["name"]
		   if !ok {
		     logger.Println("Missing repo name")
		   }
		   name, ok := iName.(string)
		   if !ok {
		     logger.Println("Invalid repo name received, got %v", iName)
		   }
		*/
	}
	repos.Store(reposArr)
	return nil
}

type ReposPageData struct {
	Repos []RepoInfo
}

func NewReposPageData() ReposPageData {
	return ReposPageData{Repos: repos.Load()}
}

type RepoInfo struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
}
