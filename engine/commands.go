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
	assignReviewerCommandRegex = regexp.MustCompile(`^\/reviewpad assign-reviewers\s+((?:[a-zA-Z0-9\-]+[a-zA-Z0-9])(?:,\s*[a-zA-Z0-9\-]+[a-zA-Z0-9])*)(?:\s+(\d+))?(?:\s+(random|round-robin|reviewpad))?$`)
)

var commands = map[*regexp.Regexp]func(matches []string) (*PadRule, *PadWorkflow, error){
	assignReviewerCommandRegex: AssignReviewerCommand,
}

func AssignReviewerCommand(matches []string) (*PadRule, *PadWorkflow, error) {
	var err error

	if len(matches) < 2 {
		return nil, nil, errors.New("invalid assign reviewer command")
	}

	reviewersList := strings.Split(strings.ReplaceAll(matches[1], " ", ""), ",")

	reviewers := `"` + strings.Join(reviewersList, `","`) + `"`

	requiredReviewers := len(reviewersList)

	if len(matches) > 2 {
		requiredReviewers, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, nil, err
		}
	}

	strategy := "reviewpad"

	if len(matches) > 3 && matches[3] != "" {
		strategy = matches[3]
	}

	name := fmt.Sprintf(`assign-reviewer-%s`, utils.RandomString(10))

	action := fmt.Sprintf(`$assignReviewer([%s], %d, %q)`, reviewers, requiredReviewers, strategy)

	rule := &PadRule{
		Name: name,
		Spec: "true",
	}

	workflow := &PadWorkflow{
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
	}

	return rule, workflow, nil
}
