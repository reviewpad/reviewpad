// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
)

var (
	assignReviewerCommandRegex = regexp.MustCompile(`^\/reviewpad assign-reviewers\s+([a-zA-Z0-9\-]+[a-zA-Z0-9],\s*[a-zA-Z0-9\-]+[a-zA-Z0-9]*)\s+(\d+)\s+(random|round-robin|reviewpad)$`)
)

var commands = map[*regexp.Regexp]func(*ReviewpadFile, []string) error{
	assignReviewerCommandRegex: assignReviewerCommand,
}

func assignReviewerCommand(file *ReviewpadFile, matches []string) error {
	if len(matches) != 4 {
		return errors.New("invalid assign reviewer command")
	}

	reviewersList := strings.Split(strings.ReplaceAll(matches[1], " ", ""), ",")

	reviewers := `"` + strings.Join(reviewersList, `","`) + `"`

	requiredReviewers, err := strconv.Atoi(matches[2])
	if err != nil {
		return err
	}

	strategy := matches[3]

	name := fmt.Sprintf(`assign-reviewer-%s`, utils.RandomString(10))

	file.Rules = append(file.Rules, PadRule{
		Name: name,
		Spec: "true",
	})

	action := fmt.Sprintf(`$assignReviewer([%s], %d, %q)`, reviewers, requiredReviewers, strategy)

	file.Workflows = append(file.Workflows, PadWorkflow{
		Name:        name,
		On:          []handler.TargetEntityKind{handler.PullRequest},
		Description: "assign reviewer command",
		Rules: []PadWorkflowRule{
			{
				Rule: name,
			},
		},
		AlwaysRun: true,
		Actions:   []string{action},
	})

	return nil
}
