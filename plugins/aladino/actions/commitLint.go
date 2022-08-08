// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-conventionalcommits"
	"github.com/reviewpad/go-conventionalcommits/parser"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	consts "github.com/reviewpad/reviewpad/v3/plugins/aladino/consts"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func CommitLint() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: commitLintCode,
	}
}

func commitLintCode(e aladino.Env, _ []aladino.Value) error {
	evalEnv := e.(*aladino.BaseEnv)
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	ghCommits, err := utils.GetPullRequestCommits(e.GetCtx(), e.GetClient(), owner, repo, prNum)
	if err != nil {
		return err
	}

	for _, ghCommit := range ghCommits {
		commitMsg := ghCommit.Commit.GetMessage()
		res, err := parser.NewMachine(conventionalcommits.WithTypes(conventionalcommits.TypesConventional)).Parse([]byte(commitMsg))

		if err != nil || !res.Ok() {
			comments := evalEnv.GetComments()
			body := fmt.Sprintf("**Unconventional commit detected**: '%v' (%v)", commitMsg, ghCommit.GetSHA())
			errors, ok := comments[consts.ERROR_LEVEL]
			if !ok {
				comments[consts.ERROR_LEVEL] = []string{body}
			} else {
				comments[consts.ERROR_LEVEL] = append(errors, body)
			}
		}
	}

	return nil
}
