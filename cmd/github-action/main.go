// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Env struct {
	RepoOwner string
	RepoName  string
	Token     string
	PRNumber  int
	Dryrun    bool
}

func getEnvVariable(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatal(fmt.Errorf("team-edition: missing %s env variable", name))
	}
	return value
}

func getEnv() *Env {
	repo := getEnvVariable("INPUT_REPOSITORY")
	token := getEnvVariable("INPUT_TOKEN")
	prnumStr := getEnvVariable("INPUT_PRNUMBER")
	dryrun := getEnvVariable("INPUT_DRYRUN") == "true"

	prNum, err := strconv.Atoi(prnumStr)
	if err != nil {
		log.Fatal(fmt.Errorf("team-edition: retrieving pull request number %v", err))
	}

	repoStrings := strings.Split(repo, "/")
	repoOwner := repoStrings[0]
	repoName := repoStrings[1]

	return &Env{
		RepoOwner: repoOwner,
		RepoName:  repoName,
		Token:     token,
		PRNumber:  prNum,
		Dryrun:    dryrun,
	}
}

func main() {
	env := getEnv()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	clientGQL := githubv4.NewClient(tc)

	ghPullRequest, _, err := client.PullRequests.Get(ctx, env.RepoOwner, env.RepoName, env.PRNumber)
	if err != nil {
		log.Fatal(err)
	}

	headRepoOwner := *ghPullRequest.Head.Repo.Owner.Login
	headRepoName := *ghPullRequest.Head.Repo.Name
	headRef := *ghPullRequest.Head.Ref

	if *ghPullRequest.Merged {
		if ghPullRequest.Head == nil {
			log.Printf("team-edition: pull request is merged and head branched was automatically deleted\n")
			return
		}

		if ghPullRequest.Head.Repo == nil {
			log.Printf("team-edition: pull request is merged and head repo branched was automatically deleted\n")
			return
		}

		if ghPullRequest.Head.Repo.DeleteBranchOnMerge == nil || *ghPullRequest.Head.Repo.DeleteBranchOnMerge {
			log.Printf("team-edition: pull request is merged and head branched was automatically deleted\n")
			return
		}

		_, _, err := client.Repositories.GetBranch(ctx, headRepoOwner, headRepoName, headRef, true)
		if err != nil {
			log.Printf("team-edition: pull request is merged and head branched cannot be retrieved\n")
			return
		}
	}

	ioReader, _, err := client.Repositories.DownloadContents(ctx, headRepoOwner, headRepoName, "reviewpad.yml", &github.RepositoryContentGetOptions{
		Ref: headRef,
	})
	if err != nil {
		log.Print(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(ioReader)

	err = reviewpad.Run(ctx, client, clientGQL, ghPullRequest, buf, env.Dryrun)

	if err != nil {
		log.Print(err.Error())
		return
	}
}
