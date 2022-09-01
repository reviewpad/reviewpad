// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"log"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func RemoveLabels() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code:           removeLabelsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func removeLabelsCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget()
	entity := t.GetTargetEntity()

	aladinoLabels := args[0].(*aladino.ArrayValue).Vals

	if len(aladinoLabels) == 0 {
		return fmt.Errorf("removeLabels: no labels provided")
	}

	repoLabels, _, err := e.GetGithubClient().GetRepositoryLabels(e.GetCtx(), entity.Owner, entity.Repo)
	if err != nil {
		return err
	}

	for _, aladinoLabel := range aladinoLabels {
		label := aladinoLabel.(*aladino.StringValue).Val

		if !labelExistsInRepo(label, repoLabels) {
			log.Printf("[warn]: the \"%v\" label does not exist in the repository", label)
			continue
		}

		labels, err := t.GetLabels()
		if err != nil {
			return err
		}

		if isLabelAppliedToIssue(label, labels) {
			err := t.RemoveLabel(label)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[warn]: the \"%v\" label is not applied to the issue", label)
		}
	}

	return nil
}

func labelExistsInRepo(label string, repoLabels []*github.Label) bool {
	for _, repoLabel := range repoLabels {
		if label == repoLabel.GetName() {
			return true
		}
	}

	return false
}

func isLabelAppliedToIssue(label string, issuesLabels []*codehost.Label) bool {
	for _, issueLabel := range issuesLabels {
		if issueLabel.Name == label {
			return true
		}
	}

	return false
}
