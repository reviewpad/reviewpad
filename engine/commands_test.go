// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
)

func TestAssignReviewer(t *testing.T) {
	testCases := map[string]struct {
		matches     []string
		wantRule    *engine.PadRule
		wantActions []string
		wantErr     error
	}{
		"when arguments are empty": {
			matches: []string{},
			wantErr: errors.New("invalid assign reviewer command"),
		},
		"when invalid number of arguments": {
			matches: []string{
				"/reviewpad assign-reviewers",
			},
			wantErr: errors.New("invalid assign reviewer command"),
		},
		"when number of reviewers is not a number": {
			matches: []string{
				"/reviewpad assign-reviewers john, jane, john2, jane27 z random",
				"john, jane, john2, jane27",
				"z",
				"random",
			},
			wantErr: &strconv.NumError{Func: "Atoi", Num: "z", Err: errors.New("invalid syntax")},
		},
		"when missing number of reviewers and policy": {
			matches: []string{
				"/reviewpad assign-reviewers john",
				"john",
				"",
				"",
			},
			wantRule: &engine.PadRule{
				Spec: "true",
			},
			wantActions: []string{`$assignReviewer(["john"], 1, "reviewpad")`},
		},
		"when missing policy": {
			matches: []string{
				"/reviewpad assign-reviewer john-123, jane 1",
				"john-123, jane",
				"1",
				"",
			},
			wantRule: &engine.PadRule{
				Spec: "true",
			},
			wantActions: []string{`$assignReviewer(["john-123","jane"], 1, "reviewpad")`},
		},
		"when only one reviewer is provided.": {
			matches: []string{
				"/reviewpad assign-reviewer john-123-jane 1 reviewpad",
				"john-123-jane",
				"1",
				"reviewpad",
			},
			wantRule: &engine.PadRule{
				Spec: "true",
			},
			wantActions: []string{`$assignReviewer(["john-123-jane"], 1, "reviewpad")`},
		},
		"when only two reviewers is provided.": {
			matches: []string{
				"/reviewpad assign-reviewer jane, john 2 random",
				"jane, john",
				"2",
				"random",
			},
			wantRule: &engine.PadRule{
				Spec: "true",
			},
			wantActions: []string{`$assignReviewer(["jane","john"], 2, "random")`},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actions, err := engine.AssignReviewerCommand(test.matches)

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantActions, actions)
		})
	}
}
