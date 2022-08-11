// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Rebase() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: rebaseCode,
	}
}

func rebaseCode(e aladino.Env, args []aladino.Value) error {
	baseOwner := gh.GetPullRequestBaseOwnerName(e.GetPullRequest())
	baseRepo := gh.GetPullRequestBaseRepoName(e.GetPullRequest())
	number := gh.GetPullRequestNumber(e.GetPullRequest())

	// Check linear history
	// https://github.com/go-git/go-git/issues/103

	// TODO: Provide pull request within the environment.
	// e.pr should be a pointer to the pull request fetched with GetPullRequest2.
	githubClient := e.GetGithubClient()
	pr, err := githubClient.GetPullRequest2(e.GetCtx(), baseOwner, baseRepo, number)
	if err != nil {
		return err
	}

	githubToken := os.Getenv("INPUT_TOKEN")

	if !pr.Rebaseable {
		return errors.New("the pull request is not rebaseable")
	}

	// CLONE

	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	auth := &http.BasicAuth{
		Username: "test", // anything except an empty string
		Password: githubToken,
	}

	// Clones the repository into the given dir, just as a normal git clone does
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		// Question: Is the action able to clone from a fork?
		URL: pr.Head.CloneURL,
		// checkr ref
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", pr.Head.Ref)),
		Auth:          auth,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[info]: rebasing %v", repo)

	if err = os.Chdir(dir); err != nil {
		return err
	}

	cmd := exec.Command("git", "remote", "add", "upstream", pr.Base.SSHURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		log.Fatalf("failed add upstream: %w", err)
		return nil
	}

	cmd = exec.Command("git", "remote", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed list remotes: %w", err)
	}

	cmd = exec.Command("git", "fetch", "upstream", pr.Base.Ref)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase: %w", err)
	}

	// cmd = exec.Command("git", "config", "user.name", "me")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err = cmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to set user.name: %w", err)
	// }

	// cmd = exec.Command("git", "config", "user.email", "me@example.com")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err = cmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to set user.email: %w", err)
	// }

	cmd = exec.Command("git", "rebase", "upstream/"+pr.Base.Ref)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase: %w", err)
	}

	// cmd = exec.Command("git", "config")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err = cmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to rebase: %w", err)
	// }

	refSpec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", pr.Head.Ref, pr.Head.Ref))

	err = repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{refSpec},
		Force:    true,
		Auth:     auth,
	})
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	return nil

	// c, _ := repo2.Config()
	// log.Printf("config: %v", c)

	// remote, err := repo2.CreateRemote(&config.RemoteConfig{
	// 	Name: "upstream",
	// 	URLs: []string{pr.Base.CloneURL},
	// })

	// c2, _ := repo2.Config()
	// log.Printf("config: %v", remote)
	// log.Printf("config: %v", c2)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // FETCH
	// err = repo2.Fetch(&git.FetchOptions{
	// 	RemoteName: "upstream",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // CHECKOUT
	// w, err := repo2.Worktree()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = w.Checkout(&git.CheckoutOptions{
	// 	Branch: plumbing.NewBranchReferenceName(pr.Head.Ref),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // MOVE TO =s branch
	// _, err := w.

	// REBASE
	// err = w.Rebase(&git.RebaseOptions{
	// 	RemoteName: "upstream",
	// 	Upstream:   plumbing.NewBranchReferenceName(pr.Base.Ref),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Prints the content of the CHANGELOG file from the cloned repository
	// changelog, err := os.Open(filepath.Join(dir, "README.md"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// io.Copy(os.Stdout, changelog)

	return nil
}
