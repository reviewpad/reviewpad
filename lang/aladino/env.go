// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/utils"
	"github.com/shurcooL/githubv4"
)

const maxPerPage = int32(100)

type TypeEnv struct {
	Vars map[string]Type
}

type EvalEnv struct {
	Ctx         context.Context
	Client      *github.Client
	ClientGQL   *githubv4.Client
	PullRequest *github.PullRequest
	Patch       map[string]*File
	RegisterMap map[string]Value
	BuiltIns    *BuiltIns
}

func NewTypeEnv(e *EvalEnv) *TypeEnv {
	builtInsType := make(map[string]Type)
	for builtInName, builtInFunction := range e.BuiltIns.Functions {
		builtInsType[builtInName] = builtInFunction.Type
	}

	for builtInName, builtInAction := range e.BuiltIns.Actions {
		builtInsType[builtInName] = builtInAction.Type
	}

	return &TypeEnv{
		Vars: builtInsType,
	}
}

func NewEvalEnv(
	ctx context.Context,
	client *github.Client,
	clientGQL *githubv4.Client,
	pullRequest *github.PullRequest,
	builtIns *BuiltIns,
) (*EvalEnv, error) {
	owner := *pullRequest.Base.Repo.Owner.Login
	repo := *pullRequest.Base.Repo.Name
	number := *pullRequest.Number

	files, err := getPullRequestFiles(ctx, client, owner, repo, number)
	if err != nil {
		return nil, err
	}

	patchMap := make(map[string]*File)

	for _, file := range files {
		patchFile, err := NewFile(file)
		if err != nil {
			return nil, err
		}

		patchMap[file.GetFilename()] = patchFile
	}

	input := &EvalEnv{
		Ctx:         ctx,
		Client:      client,
		ClientGQL:   clientGQL,
		PullRequest: pullRequest,
		Patch:       patchMap,
		RegisterMap: make(map[string]Value),
		BuiltIns:    builtIns,
	}

	return input, nil
}

func getPullRequestFiles(ctx context.Context, client *github.Client, owner string, repo string, number int) ([]*github.CommitFile, error) {
	fs, err := utils.PaginatedRequest(
		func() interface{} {
			return []*github.CommitFile{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			fls := i.([]*github.CommitFile)
			fs, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: int(maxPerPage),
			})
			if err != nil {
				return nil, nil, err
			}
			fls = append(fls, fs...)
			return fls, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return fs.([]*github.CommitFile), nil
}
